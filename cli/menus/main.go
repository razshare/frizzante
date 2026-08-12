package menus

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/razshare/frizzante/v2/cli/actions"
	"github.com/razshare/frizzante/v2/cli/apps"
	"github.com/razshare/frizzante/v2/cli/generations"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/tui/configs"
	"github.com/razshare/frizzante/v2/tui/confirm"
	"github.com/razshare/frizzante/v2/tui/search"
)

func init() {
	if files.IsFile("ask.md") {
		Main.Items = append(Main.Items, Item{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"ask", "ask.md"}, value)
			},
			Choice: search.Choice{Id: "ask", Description: "ask your local ollama (optional and experimental)"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				var askMdData []byte
				if askMdData, err = os.ReadFile("ask.md"); err != nil {
					return
				}
				if !files.IsFile(filepath.Join(".gen", "ask.docs.md")) {
					var yes bool
					if yes, err = confirm.Send(true, "ask documentation for frizzante seems to be missing; generate?"); err != nil {
						return
					}
					if yes {
						if err = generations.AskDocs(generations.AskDocsOptions{}); err != nil {
							return
						}
					}
				}
				var getAskDocsMdData []byte
				if getAskDocsMdData, err = os.ReadFile(filepath.Join(".gen", "ask.docs.md")); err != nil {
					return
				}
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ ask (optional and experimental)"))
				err = actions.Ask(actions.AskOptions{
					Model:      *app.AskModel,
					Protocol:   *app.AskProtocol,
					Host:       *app.AskHost,
					Go:         *app.Go,
					Bun:        *app.Bun,
					Tags:       *app.Tags,
					Output:     *app.Output,
					UserPrompt: strings.TrimSpace(string(askMdData)),
					SystemPrompt: strings.TrimSpace(
						strings.Join(
							[]string{
								"# Instructions",
								"1) You are an expert assistant developer that uses Frizzante to create web applications.",
								"2) You are very eager to write new code, modify existing code, read existing code, you want to help really bad.",
								"3) Your answers are being viewed in a terminal window, which could be small, so keep the width of your message lines to a maximum of 50 characters and don't use characters that make it difficult to read in the terminal, like inverted quotes (``) and things of that nature.",
								"4) Always read the files of the project to get an understanding (except for `node_modules` directories), dig deeper into the file system until you read all files (except directories like `node_modules`).",
								"5) Use your tools to look it up and learn about it.",
								"6) This whole message will be automatically injected into the conversation between you and the user every 10 messages in order to remind you what your job is.",
								"7) Don't waste words, be direct when replying to the user.",
								"8) Don't ever ask for confirmation, just do the changes.",
								"",
								"",
								"# Documentation for Frizzante",
								string(getAskDocsMdData),
							},
							"\n",
						),
					),
				})
				return
			},
		})
	}
	Main.Items = append(Main.Items, Item{
		Hidden: true,
		Choice: search.Choice{Id: "render main menu"},
		Active: func(menu *Menu, app apps.App, value string, query []string) bool { return true },
		Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
			if *app.Strict {
				err = actions.Help(actions.HelpOptions{})
				return
			}
			var data []byte
			if data, err = app.Efs.ReadFile("logo.txt"); err != nil {
				return
			}
			fmt.Println(configs.Styles.BigText.PaddingLeft(1).PaddingRight(1).Render(string(data)))
			for {
				if _, err = Render(menu, app, value, query, depth+1); err != nil {
					return
				}
				if depth > 1 {
					return
				}
			}
		},
	})
}

var Main = Menu{
	Title: "main",
	Items: []Item{
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"create", "c"}, value)
			},
			Choice: search.Choice{Id: "create project", Description: "creates a new project"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ create project"))
				err = actions.CreateProject(actions.CreateProjectOptions{
					Strict: *app.Strict,
					Efs:    app.Efs,
					Name:   value,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"dev", "d"}, value)
			},
			Choice: search.Choice{Id: "dev", Description: "runs air and vite in parallel"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ dev"))
				err = actions.Dev(actions.DevOptions{
					Go:  *app.Go,
					Air: *app.Air,
					Bun: *app.Bun,
					Efs: app.Efs,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "configure" },
			Choice: search.Choice{Id: "configure", Description: "generates binaries, installs packages and creates app/dist"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ configure"))
				err = actions.Configure(actions.ConfigureOptions{
					Go:     *app.Go,
					Air:    *app.Air,
					Bun:    *app.Bun,
					Tags:   *app.Tags,
					Output: *app.Output,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"install", "i"}, value)
			},
			Choice: search.Choice{Id: "install", Description: "installs packages"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ install"))
				err = actions.Install(actions.InstallOptions{
					Go:  *app.Go,
					Bun: *app.Bun,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"update", "u"}, value)
			},
			Choice: search.Choice{Id: "update", Description: "updates packages"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ update"))
				err = actions.Update(actions.UpdateOptions{
					Go:  *app.Go,
					Bun: *app.Bun,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "add" },
			Choice: search.Choice{Id: "add", Description: "adds packages"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ add"))
				if *app.Strict {
					err = errors.New("adding packages in strict mode is not allowed")
					return
				}
				err = actions.Npm(actions.NpmOptions{Bun: *app.Bun})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"migrate", "m"}, value)
			},
			Choice: search.Choice{Id: "migrate", Description: "migrate project database"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ migrate"))
				err = actions.Migrate(actions.MigrateOptions{
					Go:   *app.Go,
					Tags: *app.Tags,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"build", "b"}, value)
			},
			Choice: search.Choice{Id: "build", Description: "builds project"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ build"))
				err = actions.Build(actions.BuildOptions{
					Go:     *app.Go,
					Bun:    *app.Bun,
					Tags:   *app.Tags,
					Output: *app.Output,
					Value:  value,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"prebuild"}, value)
			},
			Choice: search.Choice{Id: "prebuild", Description: "executes the ./pre project"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ prebuild"))
				err = actions.PreBuild(actions.PreBuildOptions{
					Go:   *app.Go,
					Tags: *app.Tags,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"postbuild"}, value)
			},
			Choice: search.Choice{Id: "postbuild", Description: "executes the ./post project"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ postbuild"))
				err = actions.PostBuild(actions.PostBuildOptions{
					Go:   *app.Go,
					Tags: *app.Tags,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"package", "p"}, value)
			},
			Choice: search.Choice{Id: "package", Description: "packages app"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ package"))
				err = actions.Package(actions.PackageOptions{
					Go:  *app.Go,
					Bun: *app.Bun,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "package-watch" },
			Choice: search.Choice{Id: "package-watch", Description: "watches app and packages it"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ package-watch"))
				err = actions.PackageWatch(actions.PackageWatchOptions{Bun: *app.Bun})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"generate", "g"}, value)
			},
			Choice: search.Choice{Id: "generate", Description: "generates code and resources"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				_, err = Activate(&Generate, app, append([]string{value}, query...), depth+1)
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "asm" },
			Choice: search.Choice{Id: "assembly explorer", Description: "starts the assembly explorer"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ assembly explorer"))
				err = actions.AssemblyExplorer(actions.AssemblyExplorerOptions{
					Go:   *app.Go,
					Bun:  *app.Bun,
					Tags: *app.Tags,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "check" },
			Choice: search.Choice{Id: "check", Description: "checks for code errors"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ check"))
				err = actions.Check(actions.CheckOptions{
					Go:          *app.Go,
					Bun:         *app.Bun,
					Incremental: *app.Incremental,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"format", "f"}, value)
			},
			Choice: search.Choice{Id: "format", Description: "formats svelte and go code"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ format"))
				return actions.Format(actions.FormatOptions{
					Go:  *app.Go,
					Bun: *app.Bun,
				})
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "clean" },
			Choice: search.Choice{Id: "clean", Description: "deletes .gen, .vite, app/{dist,node_modules}"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ clean"))
				err = actions.CleanProject(actions.CleanProjectOptions{Go: *app.Go})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "reset" },
			Choice: search.Choice{Id: "reset", Description: "deletes global cache"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ reset"))
				err = actions.Reset(actions.ResetOptions{})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "clear" },
			Choice: search.Choice{Id: "clear", Description: "clears terminal screen"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ clear"))
				err = actions.Clear(actions.ClearOptions{})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "lock-packages" },
			Choice: search.Choice{Id: "lock packages", Description: "locks js packages to the current exact version"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ lock packages"))
				err = actions.LockPackages(actions.LockPackagesOptions{})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "test" },
			Choice: search.Choice{Id: "test", Description: "runts tests"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ test"))
				err = actions.Test(actions.TestOptions{
					Go:  *app.Go,
					Bun: *app.Bun,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"version", "v"}, value)
			},
			Choice: search.Choice{Id: "version", Description: "shows binary version"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ version"))
				err = actions.Version(actions.VersionOptions{Efs: app.Efs})
				return
			},
		},
	},
}
