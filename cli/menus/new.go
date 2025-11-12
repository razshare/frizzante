package menus

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/razshare/frizzante/cli/actions"
	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/cli/paths"
	tags_ "github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/platforms"
	"github.com/razshare/frizzante/tui/configs"
	"github.com/razshare/frizzante/tui/inputs"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
)

func New(app *apps.App) (*Menu, error) {
	cache, err := paths.Cache()
	if err != nil {
		return nil, err
	}

	platform := platforms.Detect()

	_go, err := paths.Go(*app.Go)
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
				Handler: func() error {
					return actions.Configure(actions.ConfigureOptions{
						Auto:     *app.Yes,
						Platform: platform,
						Go:       _go,
						Air:      air,
						Bun:      bun,
						Efs:      app.Efs,
					})
				},
			},
			{
				Choice: search.Choice{Id: "create project", Description: "creates a new project"},
				Active: func() bool { return *app.CreateProject != "" },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ create project (creates a new project)")))
					return actions.CreateProject(actions.CreateProjectOptions{
						Name: *app.CreateProject,
						Go:   _go,
						Efs:  app.Efs,
						Air:  air,
						Bun:  bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "install", Description: "installs go and js packages"},
				Active: func() bool { return *app.Install },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ install (installs go and js packages)")))
					return actions.Install(actions.InstallOptions{
						Go:  _go,
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
						Go:  _go,
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
					var tags []string
					if tags, err = tags_.Parse(*app.Tags); err != nil {
						return
					}

					if !*app.Dev {
						if tags, err = tags_.Select([]search.Choice{
							{Id: "types", Description: "enables type generations"},
							{Id: "no_js_runtime", Description: "disables the server-side JavaScript runtime"},
							{Id: "experimental_qjs_runtime", Description: "replaces goja with qjs"},
							{Id: "other", Description: "adds custom tags"},
						}); err != nil {
							return
						}
					}

					tags = append(tags, "dev", "trace")

					err = actions.Dev(actions.DevOptions{
						Go:   _go,
						Air:  air,
						Bun:  bun,
						Tags: tags,
						Efs:  app.Efs,
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
					var tags []string
					if tags, err = tags_.Parse(*app.Tags); err != nil {
						return
					}

					if !*app.Build {
						if tags, err = tags_.Select([]search.Choice{
							{Id: "trace", Description: "enables tracing with stack.Trace()"},
							{Id: "no_js_runtime", Description: "disables the server-side JavaScript runtime"},
							{Id: "experimental_qjs_runtime", Description: "replaces goja with qjs"},
							{Id: "other", Description: "adds custom tags"},
						}); err != nil {
							return
						}
					}

					err = actions.Build(actions.BuildOptions{
						Platform: platform,
						Go:       _go,
						Bun:      bun,
						Tags:     tags,
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
					var tags []string
					if tags, err = tags_.Parse(*app.Tags); err != nil {
						return
					}

					err = actions.AssemblyExplorer(actions.AssemblyExplorerOptions{
						Platform: platform,
						Go:       _go,
						Bun:      bun,
						Tags:     tags,
						Auto:     *app.Yes,
					})

					return
				},
			},
			{
				Choice: search.Choice{Id: "generate", Description: "generates code and resources"},
				Active: func() bool { return *app.Generate != "" },
				Handler: func() (err error) {
					var tags []string
					tags, err = tags_.Parse(*app.Tags)
					tags = append(tags, "dev")

					generation := *app.Generate

					if generation == ":pick" {
						generation = ""
					}

					err = actions.Generate(actions.GenerateOptions{
						Generation: generation,
						Auto:       *app.Yes,
						Efs:        app.Efs,
						Platform:   platform,
						Go:         _go,
						Air:        air,
						Bun:        bun,
						Sqlc:       sqlc,
						Tags:       tags,
						SqlcYaml:   *app.SqlcYaml,
						Database:   *app.Database,
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
					var offset string
					var target string

					parts := strings.SplitN(*app.Migrate, ",", 2)

					if len(parts) >= 1 {
						offset = parts[0]
					} else {
						offset = ""
					}

					if len(parts) >= 2 {
						target = parts[1]
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

					err = actions.Migrate(actions.MigrateOptions{
						Auto:     *app.Yes,
						Platform: platform,
						Sqlc:     sqlc,
						SqlcYaml: *app.SqlcYaml,
						Offset:   offset,
						Target:   target,
						Database: database,
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
						Go:  _go,
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
						Go: _go,
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
				Choice: search.Choice{Id: "test", Description: "runs tests"},
				Active: func() bool { return *app.Test },
				Handler: func() error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render(fmt.Sprint("running ▷ test (runs tests)")))
					return actions.Test(actions.TestOptions{
						Go:  _go,
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
