package menus

import (
	"errors"
	"fmt"

	"github.com/razshare/frizzante/cli/action"
	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/cli/path"
	tags_ "github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/cli/user"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
)

func New(app *apps.App) (*Menu, error) {
	cache, err := user.FrizzanteCache()
	if err != nil {
		return nil, err
	}

	plat, err := user.Platform(app)
	if err != nil {
		return nil, err
	}

	_go, err := path.Go(*app.Go)
	if err != nil {
		return nil, err
	}

	air, err := path.Air(*app.Air)
	if err != nil {
		return nil, err
	}

	bun, err := path.Bun(*app.Bun)
	if err != nil {
		return nil, err
	}

	sqlc, err := path.Sqlc(*app.Sqlc)
	if err != nil {
		return nil, err
	}

	return &Menu{
		Items: []Item{
			{
				Choice: search.Choice{Id: "configure", Description: "installs required binaries and packages"},
				Active: func() bool { return *app.Configure },
				Handler: func() error {
					return action.Configure(action.ConfigureOptions{
						App:      *app.App,
						Auto:     *app.Yes,
						Platform: plat,
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
					return action.CreateProject(action.CreateProjectOptions{
						Name: *app.CreateProject,
						Go:   _go,
						Efs:  app.Efs,
					})
				},
			},
			{
				Choice: search.Choice{Id: "install", Description: "installs go and js packages"},
				Active: func() bool { return *app.Install },
				Handler: func() error {
					return action.Install(action.InstallOptions{
						App: *app.App,
						Go:  _go,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "update", Description: "updates go and js packages"},
				Active: func() bool { return *app.Update },
				Handler: func() error {
					return action.Update(action.UpdateOptions{
						App: *app.App,
						Go:  _go,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "add", Description: "adds packages"},
				Active: func() bool { return *app.Add != "" },
				Handler: func() error {
					var packageType string
					packageType, err = singleselect.Send(
						[]search.Choice{
							{Id: "js", Description: fmt.Sprintf("installs js packages in %s/node_modules", *app.App)},
							//{Id: "go", Description: "installs go packages"},
						},
						"type of packages",
					)

					if err != nil {
						return err
					}

					if packageType == "js" {
						return action.Npm(action.NpmOptions{
							App: *app.App,
							Bun: bun,
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
					var tags []string
					if tags, err = tags_.Parse(*app.Tags); err != nil {
						return
					}

					if !*app.Dev {
						if tags, err = tags_.Select([]search.Choice{
							{Id: "no_js_runtime", Description: "disables the server-side JavaScript runtime"},
							{Id: "experimental_qjs_runtime", Description: "replaces goja with qjs"},
							{Id: "other", Description: "adds custom tags"},
						}); err != nil {
							return
						}
					}

					tags = append(tags, "dev", "types", "trace")

					err = action.Dev(action.DevOptions{
						App:  *app.App,
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

					err = action.Build(action.BuildOptions{
						App:      *app.App,
						Platform: plat,
						Go:       _go,
						Bun:      bun,
						Tags:     tags,
					})

					return
				},
			},
			{
				Choice: search.Choice{Id: "migrate", Description: "migrates database schema"},
				Active: func() bool { return *app.Migrate },
				Handler: func() error {
					return action.Migrate(action.MigrateOptions{
						Auto:     *app.Yes,
						Platform: plat,
						Sqlc:     sqlc,
						SqlcYaml: *app.SqlcYaml,
						Index:    *app.MigrateIndex,
					})
				},
			},
			{
				Choice: search.Choice{Id: "assembly explorer", Description: "explores application assembly output"},
				Active: func() bool { return *app.AssemblyExplorer },
				Handler: func() (err error) {
					var tags []string
					if tags, err = tags_.Parse(*app.Tags); err != nil {
						return
					}

					err = action.AssemblyExplorer(action.AssemblyExplorerOptions{
						App:      *app.App,
						Platform: plat,
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

					var selected string

					if *app.Generate != ":pick" && *app.Generate != "pick" {
						selected = *app.Generate
					}

					err = action.Generate(action.GenerateOptions{
						App:      *app.App,
						Selected: selected,
						Auto:     *app.Yes,
						Efs:      app.Efs,
						Platform: plat,
						Go:       _go,
						Air:      air,
						Bun:      bun,
						Sqlc:     sqlc,
						Tags:     tags,
						Active:   selected != "",
						SqlcYaml: *app.SqlcYaml,
					})

					return
				},
			},
			{
				Choice: search.Choice{Id: "package", Description: "builds app"},
				Active: func() bool { return *app.Package },
				Handler: func() error {
					return action.Package(action.PackageOptions{
						App: *app.App,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "package (watch)", Description: "builds app on change"},
				Active: func() bool { return *app.PackageWatch },
				Handler: func() error {
					return action.PackageWatch(action.PackageWatchOptions{
						App: *app.App,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "check", Description: "checks for code errors"},
				Active: func() bool { return *app.Check },
				Handler: func() error {
					return action.Check(action.CheckOptions{
						App: *app.App,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "format", Description: "format code"},
				Active: func() bool { return *app.Format },
				Handler: func() error {
					return action.Format(action.FormatOptions{
						App: *app.App,
						Go:  _go,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "touch", Description: "adds placeholders in app/dist"},
				Active: func() bool { return *app.Touch },
				Handler: func() error {
					return action.Touch(action.TouchOptions{
						App: *app.App,
					})
				},
			},
			{
				Choice: search.Choice{Id: "clean project", Description: "deletes .gen, .vite, app/{dist,node_modules}"},
				Active: func() bool { return *app.CleanProject },
				Handler: func() error {
					return action.CleanProject(action.CleanProjectOptions{
						App: *app.App,
						Go:  _go,
					})
				},
			},
			{
				Choice: search.Choice{Id: "reset", Description: "deletes " + cache},
				Active: func() bool { return *app.Reset },
				Handler: func() error {
					return action.Reset(action.ResetOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "clear", Description: "clears screen"},
				Active: func() bool { return *app.Clear },
				Handler: func() error {
					return action.Clear(action.ClearOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "test", Description: "runs tests"},
				Active: func() bool { return *app.Test },
				Handler: func() error {
					return action.Test(action.TestOptions{
						App: *app.App,
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
					return action.Welcome(action.WelcomeOptions{})
				},
			},
			{
				Hidden: true,
				Choice: search.Choice{Id: "help", Description: "shows the help menu"},
				Active: func() bool { return *app.Help },
				Handler: func() error {
					return action.Help(action.HelpOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "version", Description: "shows binary version"},
				Active: func() bool { return *app.Version },
				Handler: func() error {
					return action.Version(action.VersionOptions{
						Efs: app.Efs,
					})
				},
			},
		},
	}, nil
}
