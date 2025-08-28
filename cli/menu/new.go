package menu

import (
	"fmt"
	"github.com/razshare/frizzante/cli/action"
	"github.com/razshare/frizzante/cli/app"
	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/user"
	"github.com/razshare/frizzante/tui/input"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
)

func New(a *app.App) (*Menu, error) {
	home, err := user.FrizzanteHome()
	if err != nil {
		return nil, err
	}

	plat, err := user.Platform(a)
	if err != nil {
		return nil, err
	}

	_go, err := path.Go(*a.Go)
	if err != nil {
		return nil, err
	}

	air, err := path.Air(*a.Air)
	if err != nil {
		return nil, err
	}

	bun, err := path.Bun(*a.Bun)
	if err != nil {
		return nil, err
	}

	sqlc, err := path.Sqlc(*a.Sqlc)
	if err != nil {
		return nil, err
	}

	return &Menu{
		Items: []Item{
			{
				Choice: search.Choice{Id: "configure", Description: "installs required binaries and packages"},
				Active: func() bool { return *a.Configure },
				Handler: func() error {
					return action.Configure(action.ConfigureOptions{
						App:      *a.App,
						Auto:     *a.Yes,
						Platform: plat,
						Go:       _go,
						Air:      air,
						Bun:      bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "create project", Description: "creates a new project"},
				Active: func() bool { return *a.CreateProject != "" },
				Handler: func() error {
					if *a.CreateProject == "" {
						*a.CreateProject, err = input.Send("give the project a name")
						if err != nil {
							return err
						}
					}

					err = generate.Init(a.Efs)
					if err != nil {
						return err
					}

					return generate.Project(generate.ProjectOptions{
						Name: *a.CreateProject,
						Auto: *a.Yes,
						Efs:  a.Efs,
					})
				},
			},
			{
				Choice: search.Choice{Id: "install", Description: "installs dependencies"},
				Active: func() bool { return *a.Install },
				Handler: func() error {
					return action.Install(action.InstallOptions{
						App: *a.App,
						Go:  _go,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "update", Description: "updates dependencies"},
				Active: func() bool { return *a.Update },
				Handler: func() error {
					return action.Update(action.UpdateOptions{
						App: *a.App,
						Go:  _go,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "add", Description: "adds packages"},
				Active: func() bool { return *a.Add != "" },
				Handler: func() error {
					var t string
					t, err = singleselect.Send(
						[]search.Choice{
							{Id: "js", Description: fmt.Sprintf("installs js packages in %s/node_modules", *a.App)},
							//{Id: "go", Description: "installs go packages"},
						},
						"type of packages",
					)

					if err != nil {
						return err
					}

					if t == "js" {
						return action.Npm(action.NpmOptions{
							App: *a.App,
							Bun: bun,
						})
					}

					return fmt.Errorf("%s packages are not supported", t)
				},
			},
			{
				Choice: search.Choice{Id: "dev", Description: "runs air and vite in parallel"},
				Active: func() bool { return *a.Dev },
				Handler: func() error {
					return action.Dev(action.DevOptions{
						App: *a.App,
						Go:  _go,
						Air: air,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "build", Description: "builds project"},
				Active: func() bool { return *a.Build },
				Handler: func() error {
					return action.Build(action.BuildOptions{
						App:      *a.App,
						Platform: plat,
						Go:       _go,
						Bun:      bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "generate", Description: "generates code and resources"},
				Active: func() bool { return *a.Generate != "" },
				Handler: func() error {
					return action.Generate(action.GenerateOptions{
						App:      *a.App,
						Selected: *a.Generate,
						Auto:     *a.Yes,
						Efs:      a.Efs,
						Platform: plat,
						Go:       _go,
						Air:      air,
						Bun:      bun,
						Sqlc:     sqlc,
					})
				},
			},
			{
				Choice: search.Choice{Id: "package", Description: "builds app"},
				Active: func() bool { return *a.Package },
				Handler: func() error {
					return action.Pkg(action.PkgOptions{
						App: *a.App,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "package (watch)", Description: "builds app on change"},
				Active: func() bool { return *a.PackageWatch },
				Handler: func() error {
					return action.PkgWatch(action.PkgWatchOptions{
						App: *a.App,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "check", Description: "checks for code errors"},
				Active: func() bool { return *a.Check },
				Handler: func() error {
					return action.Check(action.CheckOptions{
						App: *a.App,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "format", Description: "format code"},
				Active: func() bool { return *a.Format },
				Handler: func() error {
					return action.Format(action.FormatOptions{
						App: *a.App,
						Go:  _go,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "touch", Description: "adds placeholders in app/dist"},
				Active: func() bool { return *a.Touch },
				Handler: func() error {
					return action.Touch(action.TouchOptions{
						App: *a.App,
					})
				},
			},
			{
				Choice: search.Choice{Id: "clean project", Description: "deletes .gen, .vite, app/{dist,node_modules}"},
				Active: func() bool { return *a.CleanProject },
				Handler: func() error {
					return action.CleanProject(action.CleanProjectOptions{
						App: *a.App,
						Go:  _go,
					})
				},
			},
			{
				Choice: search.Choice{Id: "reset", Description: "deletes " + home},
				Active: func() bool { return *a.Reset },
				Handler: func() error {
					return action.Reset(action.ResetOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "clear", Description: "clears screen"},
				Active: func() bool { return *a.Clear },
				Handler: func() error {
					return action.Clear(action.ClearOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "test", Description: "runs tests"},
				Active: func() bool { return *a.Test },
				Handler: func() error {
					return action.Test(action.TestOptions{
						App: *a.App,
						Go:  _go,
						Bun: bun,
					})
				},
			},
			{
				Hidden: true,
				Choice: search.Choice{Id: "welcome", Description: "shows a welcome message"},
				Active: func() bool { return *a.Welcome },
				Handler: func() error {
					return action.Welcome(action.WelcomeOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "help", Description: "shows the help menu"},
				Active: func() bool { return *a.Help },
				Handler: func() error {
					return action.Help(action.HelpOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "version", Description: "shows binary version"},
				Active: func() bool { return *a.Version },
				Handler: func() error {
					return action.Version(action.VersionOptions{
						Efs: a.Efs,
					})
				},
			},
		},
	}, nil
}
