package menus

import (
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/razshare/frizzante/cli/actions"
	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/tui/configs"
	"github.com/razshare/frizzante/tui/search"
)

var Main = Menu{
	Title: "main",
	Items: []Item{
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Create != "" },
			Choice: search.Choice{Id: "create project", Description: "creates a new project"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ create project"))
				err = actions.CreateProject(actions.CreateProjectOptions{
					Strict: *app.Strict,
					Name:   *app.Create,
					Efs:    app.Efs,
				})
				return err
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Dev },
			Choice: search.Choice{Id: "dev", Description: "runs air and vite in parallel"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ dev"))
				err = actions.Dev(actions.DevOptions{
					Go:     *app.Go,
					Air:    *app.Air,
					Bun:    *app.Bun,
					Tags:   *app.Tags,
					Efs:    app.Efs,
					Strict: *app.Strict,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Configure },
			Choice: search.Choice{Id: "configure", Description: "generates binaries, installs packages and creates app/dist"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ configure"))
				err = actions.Configure(actions.ConfigureOptions{
					Go:  *app.Go,
					Air: *app.Air,
					Bun: *app.Bun,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Install },
			Choice: search.Choice{Id: "install", Description: "installs packages"},
			Handle: func(menu *Menu, app apps.App) (err error) {
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
			Active: func(menu *Menu, app apps.App) bool { return *app.Update },
			Choice: search.Choice{Id: "update", Description: "updates packages"},
			Handle: func(menu *Menu, app apps.App) (err error) {
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
			Active: func(menu *Menu, app apps.App) bool { return *app.Add != "" },
			Choice: search.Choice{Id: "add", Description: "adds packages"},
			Handle: func(menu *Menu, app apps.App) (err error) {
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
			Active: func(menu *Menu, app apps.App) bool { return *app.Build },
			Choice: search.Choice{Id: "build", Description: "builds project"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ build"))
				err = actions.Build(actions.BuildOptions{
					Go:     *app.Go,
					Bun:    *app.Bun,
					Tags:   *app.Tags,
					Strict: *app.Strict,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate != "" },
			Choice: search.Choice{Id: "generate", Description: "generates code and resources"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				_, err = Activate(&Generate, app)
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Asm },
			Choice: search.Choice{Id: "assembly explorer", Description: "starts the assembly explorer"},
			Handle: func(menu *Menu, app apps.App) (err error) {
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
			Active: func(menu *Menu, app apps.App) bool { return *app.Migrate != "" },
			Choice: search.Choice{Id: "migrate", Description: "migrates database schema"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ migrate"))
				err = actions.Migrate(actions.MigrateOptions{
					Strict:   *app.Strict,
					SqlcYaml: *app.SqlcYaml,
					Query:    *app.Migrate,
					Database: *app.Database,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Package },
			Choice: search.Choice{Id: "package", Description: "packages the svelte application into app/dist"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ package"))
				err = actions.Package(actions.PackageOptions{Bun: *app.Bun})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.PackageWatch },
			Choice: search.Choice{Id: "package-watch", Description: "packages the svelte application when source code changes"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ package watch"))
				err = actions.PackageWatch(actions.PackageWatchOptions{Bun: *app.Bun})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Check },
			Choice: search.Choice{Id: "check", Description: "checks for code errors"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ check"))
				err = actions.Check(actions.CheckOptions{Bun: *app.Bun})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Format },
			Choice: search.Choice{Id: "format", Description: "formats svelte and go code"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ format"))
				return actions.Format(actions.FormatOptions{
					Go:  *app.Go,
					Bun: *app.Bun,
				})
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Touch },
			Choice: search.Choice{Id: "touch", Description: "adds placeholders in app/dist"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ touch"))
				return actions.Touch(actions.TouchOptions{})
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Clean },
			Choice: search.Choice{Id: "clean", Description: "deletes .gen, .vite, app/{dist,node_modules}"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ clean project"))
				err = actions.CleanProject(actions.CleanProjectOptions{Go: *app.Go})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Reset },
			Choice: search.Choice{Id: "reset", Description: "deletes global cache"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ reset"))
				err = actions.Reset(actions.ResetOptions{})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Clear },
			Choice: search.Choice{Id: "clear", Description: "clears terminal screen"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ clear"))
				err = actions.Clear(actions.ClearOptions{})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.LockJsPackages },
			Choice: search.Choice{Id: "lock packages", Description: "locks js packages to the current exact version"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ lock packages"))
				err = actions.LockPackages(actions.LockPackagesOptions{})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Test },
			Choice: search.Choice{Id: "test", Description: "runts tests"},
			Handle: func(menu *Menu, app apps.App) (err error) {
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
			Hidden: true,
			Active: func(menu *Menu, app apps.App) bool { return *app.Welcome },
			Choice: search.Choice{Id: "welcome", Description: "shows a welcome message and yields without killing the process"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ welcome"))
				err = actions.Welcome(actions.WelcomeOptions{})
				return
			},
		},
		{
			Hidden: true,
			Active: func(menu *Menu, app apps.App) bool { return *app.Help },
			Choice: search.Choice{Id: "help", Description: "shows the help menu"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ help"))
				err = actions.Help(actions.HelpOptions{App: app})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Version },
			Choice: search.Choice{Id: "version", Description: "shows binary version"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("running ▷ version"))
				err = actions.Version(actions.VersionOptions{Efs: app.Efs})
				return
			},
		},
		{
			Hidden: true,
			Choice: search.Choice{Id: "render main menu"},
			Active: func(menu *Menu, app apps.App) bool { return true },
			Handle: func(menu *Menu, app apps.App) (err error) {
				if *app.Strict {
					err = actions.Help(actions.HelpOptions{})
					return
				}
				var data []byte
				if data, err = app.Efs.ReadFile("logo.txt"); err != nil {
					return
				}
				fmt.Print(configs.Styles.BigText.PaddingLeft(1).PaddingRight(1).Render(string(data)))
				for {
					if _, err = Render(menu, app); err != nil {
						return
					}
				}
			},
		},
	},
}
