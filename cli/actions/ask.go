package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/razshare/frizzante/v2/cli/services/indexing"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/tui/inputs"
	"github.com/razshare/frizzante/v2/tui/messages"
)

func Ask(options AskOptions) (err error) {
	if !slices.Contains([]string{"http", "https"}, options.Protocol) {
		err = fmt.Errorf("unknown `%s` protocol", options.Protocol)
		return
	}
	type FunctionCall struct {
		Index     int64          `json:"index"`
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}
	type ToolCall struct {
		Type     string       `json:"type"`
		Function FunctionCall `json:"function"`
	}
	type Message struct {
		Role      string     `json:"role"`
		Content   string     `json:"content"`
		ToolCalls []ToolCall `json:"tool_calls"`
		ToolName  string     `json:"tool_name"`
	}
	type Property struct {
		Type        string `json:"type"`
		Description string `json:"description"`
	}
	type Parameters struct {
		Type       string              `json:"type"`
		Required   []string            `json:"required"`
		Properties map[string]Property `json:"properties"`
	}
	type Function struct {
		Name        string     `json:"name"`
		Description string     `json:"description"`
		Parameters  Parameters `json:"parameters"`
	}
	type Tool struct {
		Type     string   `json:"type"`
		Function Function `json:"function"`
	}
	type Conversation struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
		Think    bool      `json:"think"`
		Stream   bool      `json:"stream"`
		Tools    []Tool    `json:"tools"`
	}
	tools := []Tool{
		{
			Type: "function",
			Function: Function{
				Name: "FindCurrentProjectPath",
				Description: strings.Join(
					[]string{
						"Finds the absolute path to the current project.",
					},
					"\n",
				),
				Parameters: Parameters{},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "IsFile",
				Description: strings.Join(
					[]string{
						"Checks if a file exists.",
					},
					"\n",
				),
				Parameters: Parameters{
					Type:     "object",
					Required: []string{"fileName"},
					Properties: map[string]Property{
						"directoryName": {
							Type: "string",
							Description: strings.Join(
								[]string{
									"Absolute name of the file.",
									"It must be a sub-file of the current project.",
								},
								"\n",
							),
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "IsDirectory",
				Description: strings.Join(
					[]string{
						"Checks if a directory exists.",
					},
					"\n",
				),
				Parameters: Parameters{
					Type:     "string",
					Required: []string{"directoryName"},
					Properties: map[string]Property{
						"directoryName": {
							Type: "string",
							Description: strings.Join(
								[]string{
									"Absolute name of the directory.",
									"It must be a sub-directory of the current project.",
								},
								"\n",
							),
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "ListFilesInDirectory",
				Description: strings.Join(
					[]string{
						"Lists all files in a given directory.",
					},
					"\n",
				),
				Parameters: Parameters{
					Type:     "string",
					Required: []string{"directoryName"},
					Properties: map[string]Property{
						"directoryName": {
							Type: "string",
							Description: strings.Join(
								[]string{
									"Absolute name of the directory.",
									"It must be a sub-directory of the current project.",
								},
								"\n",
							),
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "ListAllProjectFiles",
				Description: strings.Join(
					[]string{
						"Lists all files in the project.",
					},
					"\n",
				),
				Parameters: Parameters{
					Type:       "string",
					Required:   []string{},
					Properties: map[string]Property{},
				},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "ReadFile",
				Description: strings.Join(
					[]string{
						"Reads a file.",
					},
					"\n",
				),
				Parameters: Parameters{
					Type:     "string",
					Required: []string{"fileName"},
					Properties: map[string]Property{
						"fileName": {
							Type: "string",
							Description: strings.Join(
								[]string{
									"Absolute name of the file.",
									"It must be a sub-file of the current project and it must have a known extension name.",
								},
								"\n",
							),
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "ModifyCode",
				Description: strings.Join(
					[]string{
						"Modifies code in a file.",
					},
					"\n",
				),
				Parameters: Parameters{
					Type:     "string",
					Required: []string{"fileName", "fileContent"},
					Properties: map[string]Property{
						"fileName": {
							Type: "string",
							Description: strings.Join(
								[]string{
									"Absolute name of the file.",
									"It must be a sub-file of the current project and it must have a known extension name.",
								},
								"\n",
							),
						},
						"fileContent": {
							Type: "string",
							Description: strings.Join(
								[]string{
									"File content.",
								},
								"\n",
							),
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "FixCode",
				Description: strings.Join(
					[]string{
						"Fixes code in a file.",
					},
					"\n",
				),
				Parameters: Parameters{
					Type:     "string",
					Required: []string{"fileName", "fileContent"},
					Properties: map[string]Property{
						"fileName": {
							Type: "string",
							Description: strings.Join(
								[]string{
									"Absolute name of the file.",
									"It must be a sub-file of the current project and it must have a known extension name.",
								},
								"\n",
							),
						},
						"fileContent": {
							Type: "string",
							Description: strings.Join(
								[]string{
									"File content.",
								},
								"\n",
							),
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "WriteFile",
				Description: strings.Join(
					[]string{
						"Creates (if necessary) and writes content a file.",
					},
					"\n",
				),
				Parameters: Parameters{
					Type:     "string",
					Required: []string{"fileName", "fileContent"},
					Properties: map[string]Property{
						"fileName": {
							Type: "string",
							Description: strings.Join(
								[]string{
									"Absolute name of the file.",
									"It must be a sub-file of the current project and it must have a known extension name.",
								},
								"\n",
							),
						},
						"fileContent": {
							Type: "string",
							Description: strings.Join(
								[]string{
									"File content.",
								},
								"\n",
							),
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "ReadWebPage",
				Description: strings.Join(
					[]string{
						"Reads the contents of a web page recursively.",
						"This can be used to fetch documentation for specific topics if the user can provide an address.",
						"It fetches the page title and body and then looks for other hyperlinks within the page and reads those as well.",
						"It has a hardcoded limit of 2 recursions.",
					},
					"\n",
				),
				Parameters: Parameters{
					Type:     "string",
					Required: []string{"address"},
					Properties: map[string]Property{
						"address": {
							Type: "string",
							Description: strings.Join(
								[]string{
									"Web page address.",
									"It must define the protocol, for example `http://` or `https://`.",
								},
								"\n",
							),
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "Build",
				Description: strings.Join(
					[]string{
						"Builds the main project into `{current_frizzante_project_path}/.gen/bin/start`.",
						"This binary starts the main application.",
						"This is the equivalent of running `frizzante build` from the command line.",
					},
					"\n",
				),
				Parameters: Parameters{},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "Install",
				Description: strings.Join(
					[]string{
						"Installs the defined Go and Js packages for the current frizzante project.",
						"This function does not install arbitrary packages.",
						"This is the equivalent of running `frizzante install` from the command line.",
					},
					"\n",
				),
				Parameters: Parameters{},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "Update",
				Description: strings.Join(
					[]string{
						"Updates the defined Go and Js packages for the current frizzante project.",
						"This is the equivalent of running `frizzante update` from the command line.",
					},
					"\n",
				),
				Parameters: Parameters{},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "Migrate",
				Description: strings.Join(
					[]string{
						"Migrates the development database of the current project.",
						"This is the equivalent of running `frizzante migrate` from the command line.",
					},
					"\n",
				),
				Parameters: Parameters{},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "Check",
				Description: strings.Join(
					[]string{
						"Checks for code style and type errors for the current project.",
						"This is the equivalent of running `frizzante check` from the command line.",
					},
					"\n",
				),
				Parameters: Parameters{},
			},
		},
		//{
		//	Type: "function",
		//	Function: Function{
		//		Name: "Format",
		//		Description: strings.Join(
		//			[]string{
		//				"Formats the code for the current project.",
		//				"This is the equivalent of running `frizzante format` from the command line.",
		//			},
		//			"\n",
		//		),
		//		Parameters: Parameters{},
		//	},
		//},
		{
			Type: "function",
			Function: Function{
				Name: "LockPackages",
				Description: strings.Join(
					[]string{
						"Locks the Js packages for the current project to their current version by removing any version modifiers from package.json.",
						"This is the equivalent of running `frizzante lock-packages` from the command line.",
					},
					"\n",
				),
				Parameters: Parameters{},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name: "Test",
				Description: strings.Join(
					[]string{
						"Runs tests for the current project.",
						"This is the equivalent of running `frizzante test` from the command line.",
					},
					"\n",
				),
				Parameters: Parameters{},
			},
		},
	}
	type Answer struct {
		Model    string
		Messages []Message
		Message  Message
	}
	askMessages := make([]Message, 0)
	if options.SystemPrompt != "" {
		askMessages = append(askMessages, Message{
			Role:    "system",
			Content: options.SystemPrompt,
		})
	}
	if options.UserPrompt != "" {
		askMessages = append(askMessages, Message{
			Role:    "user",
			Content: options.SystemPrompt,
		})
	}
	filesRead := make([]string, 0)
	var next func(remainingRecursions int) (err error)
	next = func(remainingRecursions int) (err error) {
		if remainingRecursions == 0 {
			return
		}
		var data []byte
		if data, err = json.Marshal(Conversation{
			Messages: askMessages,
			Model:    "qwen3",
			Think:    true,
			Tools:    tools,
		}); err != nil {
			return
		}
		var response *http.Response
		if response, err = http.Post(
			fmt.Sprintf("%s://%s/api/chat", options.Protocol, options.Host),
			"application/json",
			bytes.NewBuffer(data),
		); err != nil {
			return
		}
		defer func() {
			if cerr := response.Body.Close(); cerr != nil {
				if err == nil {
					err = cerr
				}
			}
		}()
		if data, err = io.ReadAll(response.Body); err != nil {
			return
		}
		var answer Answer
		if err = json.Unmarshal(data, &answer); err != nil {
			return
		}
		messagesLocal := make([]Message, 0)
		messagesLocal = append(messagesLocal, answer.Message)
		messagesLocal = append(messagesLocal, answer.Messages...)
		askMessages = append(askMessages, messagesLocal...)
		var askError error
		for _, message := range messagesLocal {
			for _, toolCall := range message.ToolCalls {
				failure := func(err error) {
					askMessages = append(askMessages, Message{
						Role:     "tool",
						ToolName: toolCall.Function.Name,
						Content: strings.Join(
							[]string{
								"# Failure",
								err.Error(),
							},
							"\n",
						),
					})
				}
				success := func(format string, args ...string) {
					askMessages = append(askMessages, Message{
						Role:     "tool",
						ToolName: toolCall.Function.Name,
						Content: strings.Join(
							[]string{
								"# Success",
								fmt.Sprintf(format, args),
							},
							"\n",
						),
					})

				}
				switch toolCall.Function.Name {
				case "FindCurrentProjectPath":
					var workingDirectory string
					if workingDirectory, askError = os.Getwd(); askError != nil {
						messages.Infof("failed to find current project path:%v", askError)
						failure(askError)
						continue
					}
					messages.Successf("found current project path `%s`", workingDirectory)
					success(workingDirectory)
				case "IsFile":
					var workingDirectory string
					if workingDirectory, askError = os.Getwd(); err != nil {
						messages.Infof("failed to find current working directory:%v", askError)
						failure(askError)
						continue
					}
					fileName := toolCall.Function.Arguments["fileName"].(string)
					if !strings.HasPrefix(fileName, workingDirectory) {
						askError = fmt.Errorf(
							"argument `fileName` must by a sub-file of the current project directory `%s`, received `%s` instead",
							workingDirectory,
							fileName,
						)
						messages.Infof("missing file name prefix:%v", askError)
						failure(askError)
						continue
					}
					if files.IsFile(fileName) {
						messages.Successf("file `%s` exists", fileName)
						success("true")
						continue
					}
					messages.Successf("file `%s` does not exist", fileName)
					success("false")
				case "IsDirectory":
					var workingDirectory string
					if workingDirectory, askError = os.Getwd(); err != nil {
						messages.Infof("failed to find current working directory:%v", askError)
						failure(askError)
						continue
					}
					directoryName := toolCall.Function.Arguments["directoryName"].(string)
					if !strings.HasPrefix(directoryName, workingDirectory) {
						askError = fmt.Errorf(
							"argument `directoryName` must by a sub-directory of the current project directory `%s`, received `%s` instead",
							workingDirectory,
							directoryName,
						)
						messages.Infof("missing directory name prefix:%v", askError)
						failure(askError)
						continue
					}
					if files.IsDirectory(directoryName) {
						messages.Successf("directory `%s` exists", directoryName)
						success("true")
						continue
					}
					messages.Successf("directory `%s` does not exist", directoryName)
					success("false")
				case "ListFilesInDirectory":
					var workingDirectory string
					if workingDirectory, askError = os.Getwd(); err != nil {
						messages.Infof("failed to find current working directory:%v", askError)
						failure(askError)
						continue
					}
					directoryName := toolCall.Function.Arguments["directoryName"].(string)
					if !strings.HasPrefix(directoryName, workingDirectory) {
						askError = fmt.Errorf(
							"argument `directoryName` must by a sub-directory of the current project directory `%s`, received `%s` instead",
							workingDirectory,
							directoryName,
						)
						messages.Infof("missing directory name prefix:%v", askError)
						failure(askError)
						continue
					}
					type File struct {
						Name        string
						SizeInBytes int64
						IsDirectory bool
						ModifiedAt  string
					}
					filesLocal := make([]File, 0)
					var entries []os.DirEntry
					if entries, askError = os.ReadDir(directoryName); askError != nil {
						messages.Infof("failed to list files in directory `%s`:%v", directoryName, askError)
						failure(askError)
						continue
					}
					var askErrorLocal error
					for _, entry := range entries {
						var info os.FileInfo
						if info, askErrorLocal = entry.Info(); err != nil {
							break
						} else {
							filesLocal = append(filesLocal, File{
								Name:        entry.Name(),
								SizeInBytes: info.Size(),
								IsDirectory: info.IsDir(),
								ModifiedAt:  info.ModTime().String(),
							})
						}
					}
					if askErrorLocal != nil {
						askError = askErrorLocal
						messages.Infof("failed to list files in directory `%s`:%v", directoryName, askError)
						failure(askError)
						continue
					}
					var builder strings.Builder
					for _, file := range filesLocal {
						builder.WriteString(
							fmt.Sprintf(
								"- file name `%s`, size `%d bytes`, last modified at `%s`\n",
								file.Name,
								file.SizeInBytes,
								file.ModifiedAt,
							),
						)
					}
					messages.Successf("files in directory `%s` listed", directoryName)
					success(builder.String())
				case "ListAllProjectFiles":
					type File struct {
						Name        string
						SizeInBytes int64
						IsDirectory bool
						ModifiedAt  string
					}
					var list func(directoryName string) (files []File, err error)
					list = func(directoryName string) (files []File, err error) {
						files = make([]File, 0)
						var entries []os.DirEntry
						if entries, err = os.ReadDir(directoryName); err != nil {
							return
						}
						for _, entry := range entries {
							var info os.FileInfo
							if info, err = entry.Info(); err != nil {
								return
							}
							relativeName := info.Name()
							absoluteName := filepath.Join(directoryName, relativeName)
							if info.IsDir() {
								isNodeModules := relativeName == "node_modules"
								isAppDist := strings.HasSuffix(directoryName, filepath.Join("app", "dist"))
								if isNodeModules || isAppDist {
									files = append(files, File{
										Name:        absoluteName,
										SizeInBytes: info.Size(),
										IsDirectory: true,
										ModifiedAt:  info.ModTime().String(),
									})
									continue
								}
								var filesLocal []File
								if filesLocal, err = list(absoluteName); err != nil {
									return
								}
								files = append(files, filesLocal...)
								continue
							}
							files = append(files, File{
								Name:        absoluteName,
								SizeInBytes: info.Size(),
								IsDirectory: false,
								ModifiedAt:  info.ModTime().String(),
							})
						}
						return
					}
					var workingDirectory string
					if workingDirectory, askError = os.Getwd(); err != nil {
						messages.Infof("failed to find current working directory:%v", askError)
						failure(askError)
						continue
					}
					filesLocal := make([]File, 0)
					if filesLocal, askError = list(workingDirectory); err != nil {
						messages.Infof("failed to list all files in project directory `%s`:%v", workingDirectory, askError)
						failure(askError)
						continue
					}
					var builder strings.Builder
					for _, file := range filesLocal {
						builder.WriteString(
							fmt.Sprintf(
								"- file name `%s`, size `%d bytes`, last modified at `%s`\n",
								file.Name,
								file.SizeInBytes,
								file.ModifiedAt,
							),
						)
					}
					messages.Success("project files listed")
					success(builder.String())
				case "ReadFile":
					var workingDirectory string
					if workingDirectory, askError = os.Getwd(); err != nil {
						messages.Infof("failed to find current working directory:%v", askError)
						failure(askError)
						continue
					}
					fileName := toolCall.Function.Arguments["fileName"].(string)
					if !strings.HasPrefix(fileName, workingDirectory) {
						askError = fmt.Errorf(
							"argument `fileName` must by a sub-file of the current project directory `%s`, received `%s` instead",
							workingDirectory,
							fileName,
						)
						messages.Infof("missing file name prefix:%v", askError)
						failure(askError)
						continue
					}
					if data, askError = os.ReadFile(fileName); err != nil {
						messages.Infof("failed to read file `%s`:%v", fileName, askError)
						failure(askError)
						continue
					}
					filesRead = append(filesRead, fileName)
					messages.Successf("file `%s` read", fileName)
					success(string(data))
				case "FixCode", "WriteFile":
					var workingDirectory string
					if workingDirectory, askError = os.Getwd(); err != nil {
						messages.Infof("failed to find current working directory:%v", askError)
						failure(askError)
						continue
					}
					fileName := toolCall.Function.Arguments["fileName"].(string)
					if !slices.Contains(filesRead, fileName) {
						askError = fmt.Errorf("you need to read the contents of the file `%s` before attempting to write anything to it", fileName)
						failure(askError)
						continue
					}
					fileContent := toolCall.Function.Arguments["fileContent"].(string)
					if !strings.HasPrefix(fileName, workingDirectory) {
						askError = fmt.Errorf(
							"argument `fileName` must by a sub-file of the current project directory `%s`, received `%s` instead",
							workingDirectory,
							fileName,
						)
						messages.Infof("missing file name prefix:%v", askError)
						failure(askError)
						continue
					}
					if askError = os.WriteFile(fileName, []byte(fileContent), os.ModePerm); askError != nil {
						messages.Infof("failed to write to file `%s`:%v", fileName, askError)
						failure(askError)
						continue
					}
					if askError = Check(CheckOptions{
						Go:  options.Go,
						Bun: options.Bun,
					}); askError != nil {
						messages.Infof("content written to file `%s` introduced code issues:%v", fileName, askError)
						failure(askError)
						continue
					}
					messages.Successf("content written to file `%s`", fileName)
					success("content written successfully to file `%s`", fileName)
				case "ReadWebPage":
					address := toolCall.Function.Arguments["address"].(string)
					var pages map[string]indexing.IndexedPage
					if pages, askError = indexing.Index(indexing.IndexOptions{
						Address:     address,
						Context:     context.Background(),
						StickToHost: true,
						Depth:       1,
					}); askError != nil {
						messages.Infof("failed to read web page `%s`:%v", address, askError)
						failure(askError)
						continue
					}
					var builder strings.Builder
					for addressLocal, page := range pages {
						builder.WriteString(fmt.Sprintf("<!-- Beginning of page %s -->\n", page.Title))
						builder.WriteString("<html>\n")
						builder.WriteString("    <head>\n")
						builder.WriteString("        <title>\n")
						if page.Title == "" {
							builder.WriteString(addressLocal + "\n")
						} else {
							builder.WriteString(page.Title + "\n")
						}
						builder.WriteString("        </title>\n")
						builder.WriteString("    </head>\n")
						builder.WriteString("    <body>\n")
						builder.WriteString(page.Body + "\n")
						builder.WriteString("    </body>\n")
						builder.WriteString("</html>\n")
						builder.WriteString(fmt.Sprintf("<!-- Ending of page %s -->\n\n", page.Title))
					}
					messages.Successf("web page `%s` read", address)
					success(builder.String())
				case "Build":
					if askError = Build(BuildOptions{
						Go:     options.Go,
						Tags:   options.Tags,
						Output: options.Output,
						Bun:    options.Bun,
					}); askError != nil {
						messages.Infof("failed to build project:%v", askError)
						failure(askError)
						continue
					}
					messages.Success("project build successfully")
					success("project build successfully")
				case "Install":
					if askError = Install(InstallOptions{
						Go:  options.Go,
						Bun: options.Bun,
					}); askError != nil {
						messages.Infof("failed to install project packages:%v", askError)
						failure(askError)
						continue
					}
					messages.Success("packages installed successfully")
					success("packages installed successfully")
				case "Update":
					if askError = Update(UpdateOptions{
						Go:  options.Go,
						Bun: options.Bun,
					}); askError != nil {
						messages.Infof("failed to update project packages:%v", askError)
						failure(askError)
						continue
					}
					messages.Success("packages updated successfully")
					success("packages updated successfully")
				case "Migrate":
					if askError = Migrate(MigrateOptions{
						Go:   options.Go,
						Tags: options.Tags,
					}); askError != nil {
						messages.Infof("failed to migrate development database:%v", askError)
						failure(askError)
						continue
					}
					messages.Success("database migrated successfully")
					success("database migrated successfully")
				case "Check":
					if askError = Check(CheckOptions{
						Go:  options.Go,
						Bun: options.Bun,
					}); askError != nil {
						messages.Infof("check failed:%v", askError)
						failure(askError)
						continue
					}
					messages.Success("project checked successfully")
					success("project checked successfully")
				//case "Format":
				//	if askError = Format(FormatOptions{
				//		Go:  options.Bun,
				//		Bun: options.Bun,
				//	}); askError != nil {
				//		messages.Infof("failed to format project:%v", askError)
				//		failure(askError)
				//		continue
				//	}
				//	success("project formatted successfully")
				case "LockPackages":
					if askError = LockPackages(LockPackagesOptions{}); askError != nil {
						messages.Infof("failed to lock frontend packages:%v", askError)
						failure(askError)
						continue
					}
					messages.Success("project frontend packages locked successfully")
					success("project frontend packages locked successfully")
				case "Test":
					messages.Info("running tests...")
					if askError = Test(TestOptions{
						Go:  options.Go,
						Bun: options.Bun,
					}); askError != nil {
						messages.Infof("failed to run rests:%v", askError)
						failure(askError)
						continue
					}
					messages.Success("project tested successfully")
					success("project tested successfully")
				default:
					failure(fmt.Errorf("unknown function tool %s", toolCall.Function.Name))
				}
			}
			if len(message.ToolCalls) > 0 {
				if err = next(remainingRecursions - 1); err != nil {
					return
				}
			}
		}
		var trimmedContent string
		if trimmedContent = strings.TrimSpace(answer.Message.Content); trimmedContent == "" {
			return
		}
		messages.Chat(strings.ToUpper(answer.Message.Role), trimmedContent)
		return
	}
	var threshold int
	for {
		filesRead = make([]string, 0)
		if threshold < 0 {
			threshold = 50
			askMessages = append(askMessages, Message{Role: "system", Content: options.UserPrompt})
		}
		var query string
		if query, err = inputs.Send("Query"); err != nil {
			return
		}
		askMessages = append(askMessages, Message{Role: "user", Content: query})
		messages.Info("thinking...")
		if err = next(50); err != nil {
			return
		}
		threshold--
	}
}
