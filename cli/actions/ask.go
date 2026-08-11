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
	"strings"

	"github.com/razshare/frizzante/v2/cli/services/indexing"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/tui/inputs"
	messages_ "github.com/razshare/frizzante/v2/tui/messages"
	"github.com/razshare/frizzante/v2/tui/spinners"
)

func Ask(options AskOptions) (err error) {
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
		{
			Type: "function",
			Function: Function{
				Name: "Format",
				Description: strings.Join(
					[]string{
						"Formats the code for the current project.",
						"This is the equivalent of running `frizzante format` from the command line.",
					},
					"\n",
				),
				Parameters: Parameters{},
			},
		},
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
	messages := make([]Message, 0)
	var next func(remainingRecursions int) (err error)
	next = func(remainingRecursions int) (err error) {
		if remainingRecursions == 0 {
			return
		}
		var data []byte
		if data, err = json.Marshal(Conversation{
			Messages: messages,
			Model:    "qwen3",
			Think:    true,
			Tools:    tools,
		}); err != nil {
			return
		}
		var response *http.Response
		if response, err = http.Post(
			fmt.Sprintf("%s/api/chat", options.Host),
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
		messages = append(messages, messagesLocal...)
		for _, message := range messagesLocal {
			for _, toolCall := range message.ToolCalls {
				switch toolCall.Function.Name {
				case "FindCurrentProjectPath":
					var workingDirectory string
					if workingDirectory, err = os.Getwd(); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  workingDirectory,
						ToolName: toolCall.Function.Name,
					})
				case "IsFile":
					var isFile string
					if files.IsFile(toolCall.Function.Arguments["fileName"].(string)) {
						isFile = "true"
					} else {
						isFile = "false"
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  isFile,
						ToolName: toolCall.Function.Name,
					})
				case "IsDirectory":
					var isDirectory string
					if files.IsDirectory(toolCall.Function.Arguments["directoryName"].(string)) {
						isDirectory = "true"
					} else {
						isDirectory = "false"
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  isDirectory,
						ToolName: toolCall.Function.Name,
					})
				case "ListFilesInDirectory":
					type File struct {
						Name        string
						SizeInBytes int64
						IsDirectory bool
						ModifiedAt  string
					}
					filesLocal := make([]File, 0)
					var entries []os.DirEntry
					entries, err = os.ReadDir(toolCall.Function.Arguments["directoryName"].(string))
					for _, entry := range entries {
						var info os.FileInfo
						if info, err = entry.Info(); err != nil {
							return
						}
						filesLocal = append(filesLocal, File{
							Name:        entry.Name(),
							SizeInBytes: info.Size(),
							IsDirectory: info.IsDir(),
							ModifiedAt:  info.ModTime().String(),
						})
					}
					if data, err = json.MarshalIndent(filesLocal, "", " "); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  string(data),
						ToolName: toolCall.Function.Name,
					})
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
						entries, err = os.ReadDir(directoryName)
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
					if workingDirectory, err = os.Getwd(); err != nil {
						return
					}
					filesLocal := make([]File, 0)
					if filesLocal, err = list(workingDirectory); err != nil {
						return
					}
					if data, err = json.MarshalIndent(filesLocal, "", " "); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  string(data),
						ToolName: toolCall.Function.Name,
					})
				case "ReadFile":
					if data, err = os.ReadFile(toolCall.Function.Arguments["fileName"].(string)); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  string(data),
						ToolName: toolCall.Function.Name,
					})
				case "ReadWebPage":
					var pages map[string]indexing.IndexedPage
					if pages, err = indexing.Index(indexing.IndexOptions{
						Address:     toolCall.Function.Arguments["address"].(string),
						Context:     context.Background(),
						StickToHost: true,
						Depth:       1,
					}); err != nil {
						return
					}
					if data, err = json.MarshalIndent(pages, "", " "); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  string(data),
						ToolName: toolCall.Function.Name,
					})
				case "Build":
					if err = Build(BuildOptions{
						Go:     options.Go,
						Tags:   options.Tags,
						Output: options.Output,
						Bun:    options.Bun,
					}); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  "(ok)",
						ToolName: toolCall.Function.Name,
					})
				case "Install":
					if err = Install(InstallOptions{
						Go:  options.Go,
						Bun: options.Bun,
					}); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  "(ok)",
						ToolName: toolCall.Function.Name,
					})
				case "Update":
					if err = Update(UpdateOptions{
						Go:  options.Go,
						Bun: options.Bun,
					}); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  "(ok)",
						ToolName: toolCall.Function.Name,
					})
				case "Migrate":
					if err = Migrate(MigrateOptions{
						Go:   options.Go,
						Tags: options.Tags,
					}); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  "(ok)",
						ToolName: toolCall.Function.Name,
					})
				case "Check":
					if err = Check(CheckOptions{
						Bun: options.Bun,
					}); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  "(ok)",
						ToolName: toolCall.Function.Name,
					})
				case "Format":
					if err = Format(FormatOptions{
						Go:  options.Bun,
						Bun: options.Bun,
					}); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  "(ok)",
						ToolName: toolCall.Function.Name,
					})
				case "LockPackages":
					if err = LockPackages(LockPackagesOptions{}); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  "(ok)",
						ToolName: toolCall.Function.Name,
					})
				case "Test":
					if err = Test(TestOptions{
						Go:  options.Go,
						Bun: options.Bun,
					}); err != nil {
						return
					}
					messages = append(messages, Message{
						Role:     "tool",
						Content:  "(ok)",
						ToolName: toolCall.Function.Name,
					})
				default:
					return fmt.Errorf("unknown function tool %s", toolCall.Function.Name)
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
		messages_.Chat(answer.Message.Role, trimmedContent)
		return
	}
	var threshold int
	for {
		if threshold < 0 {
			threshold = 10
			messages = append(messages, Message{Role: "system", Content: options.SystemPrompt})
		}
		var query string
		if query, err = inputs.Send("Query"); err != nil {
			return
		}
		messages = append(messages, Message{Role: "user", Content: query})
		messages_.Chat("user", query)
		spinner := spinners.New("thinking...")
		go spinners.Start(spinner)
		if err = next(5); err != nil {
			spinners.Stop(spinner)
			return
		}
		spinners.Stop(spinner)
		threshold--
	}
}
