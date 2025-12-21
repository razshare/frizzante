package menus

import (
	"database/sql"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/razshare/frizzante/cli/actions"
	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/cli/paths"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/platforms"
	"github.com/razshare/frizzante/tui/configs"
	"github.com/razshare/frizzante/tui/inputs"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
)

func New(app *apps.App) (*Menu, error) {
	cache, err := paths.Cache()
	if err != nil {
		return nil, err
	}

	platform := platforms.Detect()

	go_, err := paths.Go(*app.Go)
	if err != nil {
		return nil, err
	}

	air, err := paths.Air(*app.Air)
	if err != nil {
		return nil, err
	}

	bun, err := paths.Bun(*app.Bun)
	if err != nil {
		return nil, err
	}

	sqlc, err := paths.Sqlc(*app.Sqlc)
	if err != nil {
		return nil, err
	}

	return &Menu{
		Items: []Item{
			{
				Choice: search.Choice{Id: "configure", Description: "generates bun and air binaries"},
				Active: func() bool { return *app.Configure },
				Handler: func() (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ configure (generates bun and air binaries in .gen)")))

					if _, err = exec.LookPath(air); err != nil || !files.IsFile(air) {
						messages.Info(err)
						if err = generate.Air(generate.AirOptions{Air: air, Platform: platform}); err != nil {
							return
						}
					}

					if _, err = exec.LookPath(bun); err != nil || !files.IsFile(bun) {
						messages.Info(err)
						if err = generate.Bun(generate.BunOptions{Bun: bun, Platform: platform}); err != nil {
							return
						}
					}

					err = actions.Configure(actions.ConfigureOptions{
						Go:       go_,
						Air:      air,
						Bun:      bun,
						Efs:      app.Efs,
						Platform: platform,
					})

					return
				},
			},
			{
				Choice: search.Choice{Id: "create project", Description: "creates a new project"},
				Active: func() bool { return *app.CreateProject != "" },
				Handler: func() (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ create project (creates a new project)")))

					var name string
					if name = *app.CreateProject; name == "" {
						name, err = inputs.Send("give the project a name")
						if err != nil {
							return
						}
					}

					err = actions.CreateProject(actions.CreateProjectOptions{
						Name: name,
						Go:   go_,
						Efs:  app.Efs,
						Air:  air,
						Bun:  bun,
					})
					return err
				},
			},
			{
				Choice: search.Choice{Id: "install", Description: "installs go and js packages"},
				Active: func() bool { return *app.Install },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ install (installs go and js packages)")))
					return actions.Install(actions.InstallOptions{
						Go:  go_,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "update", Description: "updates go and js packages"},
				Active: func() bool { return *app.Update },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ update (updates go and js packages)")))
					return actions.Update(actions.UpdateOptions{
						Go:  go_,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "add", Description: "adds packages"},
				Active: func() bool { return *app.Add != "" },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ add (adds packages)")))
					var packageType string
					packageType, err = select_one.Send(
						[]search.Choice{
							{Id: "js", Description: fmt.Sprint("installs js packages in app/node_modules")},
							//{Id: "go", Description: "installs go packages"},
						},
						"type of packages",
					)

					if err != nil {
						return err
					}

					if packageType == "js" {
						return actions.Npm(actions.NpmOptions{
							Query: *app.Add,
							Bun:   bun,
						})
					}

					if packageType == "" {
						return errors.New("no package type selected")
					}

					return fmt.Errorf("%s packages are not supported", packageType)
				},
			},
			{
				Choice: search.Choice{Id: "dev", Description: "runs air and vite in parallel"},
				Active: func() bool { return *app.Dev },
				Handler: func() (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ dev (runs air and vite in parallel)")))
					// The "interactive" mode is enabled only when the menu option has been activated through the main menu.
					// Whenever the menu handler is activated without going through the main menu, for example
					// by inlining the flag directly, then we don't treat the program as "interactive".
					interactive := !*app.Dev
					err = actions.Dev(actions.DevOptions{
						Go:          go_,
						Air:         air,
						Bun:         bun,
						Tags:        *app.Tags,
						Efs:         app.Efs,
						Interactive: interactive,
					})
					return
				},
			},
			{
				Choice: search.Choice{Id: "build", Description: "builds project"},
				Active: func() bool { return *app.Build },
				Handler: func() (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ build (builds project)")))
					// The "interactive" mode is enabled only when the menu option has been activated through the main menu.
					// Whenever the menu handler is activated without going through the main menu, for example
					// by inlining the flag directly, then we don't treat the program as "interactive".
					interactive := !*app.Build
					err = actions.Build(actions.BuildOptions{
						Go:          go_,
						Bun:         bun,
						Tags:        *app.Tags,
						Interactive: interactive,
					})
					return
				},
			},
			{
				Choice: search.Choice{Id: "assembly explorer", Description: "explores application assembly output"},
				Active: func() bool { return *app.AssemblyExplorer },
				Handler: func() (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ assembly explorer (explores application assembly output)")))
					// The "interactive" mode is enabled only when the menu option has been activated through the main menu.
					// Whenever the menu handler is activated without going through the main menu, for example
					// by inlining the flag directly, then we don't treat the program as "interactive".
					interactive := !*app.Dev
					err = actions.AssemblyExplorer(actions.AssemblyExplorerOptions{
						Go:          go_,
						Bun:         bun,
						Tags:        *app.Tags,
						Interactive: interactive,
					})
					return
				},
			},
			{
				Choice: search.Choice{Id: "generate", Description: "generates code and resources"},
				Active: func() bool { return *app.Generate != "" },
				Handler: func() (err error) {
					err = actions.Generate(actions.GenerateOptions{
						Generation:   *app.Generate,
						Efs:          app.Efs,
						Go:           go_,
						Air:          air,
						Bun:          bun,
						Sqlc:         sqlc,
						Tags:         *app.Tags,
						SqlcYaml:     *app.SqlcYaml,
						Database:     *app.Database,
						Platform:     platform,
						DatabaseType: *app.DatabaseType,
					})

					return
				},
			},
			{
				Choice: search.Choice{Id: "migrate", Description: "migrates database schema"},
				Active: func() bool { return *app.Migrate != "" },
				Handler: func() (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ migrate (migrates database schema)")))

					// The "interactive" mode is enabled only when the menu option has been activated through the main menu.
					// Whenever the menu handler is activated without going through the main menu, for example
					// by inlining the flag directly, then we don't treat the program as "interactive".
					interactive := *app.Migrate != ""

					var offset string
					var target string

					migrateRange := strings.SplitN(*app.Migrate, ",", 2)

					if len(migrateRange) >= 1 {
						offset = migrateRange[0]
					} else {
						offset = ""
					}

					if len(migrateRange) >= 2 {
						target = migrateRange[1]
						if offset == "" {
							offset = "first"
						}

						if target == "" {
							target = "last"
						}
					} else {
						target = ""
					}

					var names []string
					if names, err = files.FindWithSuffix("lib", ".sqlite"); err != nil {
						return
					}

					choices := make([]search.Choice, len(names))
					for index, name := range names {
						choices[index] = search.Choice{Id: name}
					}

					choices = append(choices, search.Choice{Id: "other", Description: "use a different file"})

					var name string
					if name, err = select_one.Sendf(choices, "where's your sqlite database located?"); err != nil {
						return
					}

					if name == "other" {
						if name, err = inputs.Send("where's the file located?"); err != nil {
							return
						}
					}

					var database *sql.DB
					if database, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", name)); err != nil {
						return
					}

					if _, err = exec.LookPath(sqlc); err != nil && !files.IsFile(sqlc) {
						if err = generate.Sqlc(generate.SqlcOptions{
							Sqlc:     sqlc,
							Platform: platform,
						}); err != nil {
							return
						}
					}

					err = actions.Migrate(actions.MigrateOptions{
						Sqlc:        sqlc,
						SqlcYaml:    *app.SqlcYaml,
						Offset:      offset,
						Target:      target,
						Database:    database,
						Platform:    platform,
						Interactive: interactive,
					})
					return
				},
			},
			{
				Choice: search.Choice{Id: "package", Description: "builds app"},
				Active: func() bool { return *app.Package },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ package (builds app)")))
					return actions.Package(actions.PackageOptions{
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "package watch", Description: "builds app on change"},
				Active: func() bool { return *app.PackageWatch },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ package watch (builds app on change)")))
					return actions.PackageWatch(actions.PackageWatchOptions{
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "check", Description: "checks for code errors"},
				Active: func() bool { return *app.Check },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ check (checks for code errors)")))
					return actions.Check(actions.CheckOptions{
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "format", Description: "format code"},
				Active: func() bool { return *app.Format },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ format (format code)")))
					return actions.Format(actions.FormatOptions{
						Go:  go_,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "touch", Description: "adds placeholders in app/dist"},
				Active: func() bool { return *app.Touch },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ touch (adds placeholders in app/dist)")))
					return actions.Touch(actions.TouchOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "clean project", Description: "deletes .gen, .vite, app/{dist,node_modules}"},
				Active: func() bool { return *app.CleanProject },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ clean project (deletes .gen, .vite, app/{dist,node_modules})")))
					return actions.CleanProject(actions.CleanProjectOptions{
						Go: go_,
					})
				},
			},
			{
				Choice: search.Choice{Id: "reset", Description: "deletes " + cache},
				Active: func() bool { return *app.Reset },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprintf("running ▷ reset (deletes %s)", cache)))
					return actions.Reset(actions.ResetOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "clear", Description: "clears screen"},
				Active: func() bool { return *app.Clear },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ clear (clears screen)")))
					return actions.Clear(actions.ClearOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "lock packages", Description: "locks packages to the current exact version"},
				Active: func() bool { return *app.LockPackages },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ lock packages (locks packages to the current exact version)")))
					return actions.LockPackages(actions.LockPackagesOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "snapshot", Description: "snapshots the server state and generates static web assets"},
				Active: func() bool { return *app.Snapshot },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ snapshot (snapshots the server state and generates static web assets)")))
					return actions.LockPackages(actions.LockPackagesOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "test", Description: "runs tests"},
				Active: func() bool { return *app.Test },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ test (runs tests)")))
					return actions.Test(actions.TestOptions{
						Go:  go_,
						Bun: bun,
					})
				},
			},
			{
				Hidden: true,
				Choice: search.Choice{Id: "welcome", Description: "shows a welcome message"},
				Active: func() bool { return *app.Welcome },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ welcome (shows a welcome message)")))
					return actions.Welcome(actions.WelcomeOptions{})
				},
			},
			{
				Hidden: true,
				Choice: search.Choice{Id: "help", Description: "shows the help menu"},
				Active: func() bool { return *app.Help },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ help (shows the help menu)")))
					return actions.Help(actions.HelpOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "version", Description: "shows binary version"},
				Active: func() bool { return *app.Version },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ version (shows binary version)")))
					return actions.Version(actions.VersionOptions{
						Efs: app.Efs,
					})
				},
			},
		},
	}, nil
}
