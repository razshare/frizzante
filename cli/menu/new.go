package menu

import (
	"fmt"

	action2 "github.com/razshare/frizzante/cli/action"
	"github.com/razshare/frizzante/cli/app"
	path2 "github.com/razshare/frizzante/cli/path"
	user2 "github.com/razshare/frizzante/cli/user"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
)

func New(a *app.App) (*Menu, error) {
	cache, err := user2.FrizzanteCache()
	if err != nil {
		return nil, err
	}

	plat, err := user2.Platform(a)
	if err != nil {
		return nil, err
	}

	_go, err := path2.Go(*a.Go)
	if err != nil {
		return nil, err
	}

	air, err := path2.Air(*a.Air)
	if err != nil {
		return nil, err
	}

	bun, err := path2.Bun(*a.Bun)
	if err != nil {
		return nil, err
	}

	sqlc, err := path2.Sqlc(*a.Sqlc)
	if err != nil {
		return nil, err
	}

	return &Menu{
		Items: []Item{
			{
				Choice: search.Choice{Id: "configure", Description: "installs required binaries and packages"},
				Active: func() bool { return *a.Configure },
				Handler: func() error {
					return action2.Configure(action2.ConfigureOptions{
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
					return action2.CreateProject(action2.CreateProjectOptions{
						Name: *a.CreateProject,
						Efs:  a.Efs,
					})
				},
			},
			{
				Choice: search.Choice{Id: "install", Description: "installs dependencies"},
				Active: func() bool { return *a.Install },
				Handler: func() error {
					return action2.Install(action2.InstallOptions{
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
					return action2.Update(action2.UpdateOptions{
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
						return action2.Npm(action2.NpmOptions{
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
					return action2.Dev(action2.DevOptions{
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
					return action2.Build(action2.BuildOptions{
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
					var selected string

					if *a.Generate != ":pick" {
						selected = *a.Generate
					}

					return action2.Generate(action2.GenerateOptions{
						App:      *a.App,
						Selected: selected,
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
					return action2.Package(action2.PackageOptions{
						App: *a.App,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "package (watch)", Description: "builds app on change"},
				Active: func() bool { return *a.PackageWatch },
				Handler: func() error {
					return action2.PackageWatch(action2.PackageWatchOptions{
						App: *a.App,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "check", Description: "checks for code errors"},
				Active: func() bool { return *a.Check },
				Handler: func() error {
					return action2.Check(action2.CheckOptions{
						App: *a.App,
						Bun: bun,
					})
				},
			},
			{
				Choice: search.Choice{Id: "format", Description: "format code"},
				Active: func() bool { return *a.Format },
				Handler: func() error {
					return action2.Format(action2.FormatOptions{
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
					return action2.Touch(action2.TouchOptions{
						App: *a.App,
					})
				},
			},
			{
				Choice: search.Choice{Id: "clean project", Description: "deletes .gen, .vite, app/{dist,node_modules}"},
				Active: func() bool { return *a.CleanProject },
				Handler: func() error {
					return action2.CleanProject(action2.CleanProjectOptions{
						App: *a.App,
						Go:  _go,
					})
				},
			},
			{
				Choice: search.Choice{Id: "reset", Description: "deletes " + cache},
				Active: func() bool { return *a.Reset },
				Handler: func() error {
					return action2.Reset(action2.ResetOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "clear", Description: "clears screen"},
				Active: func() bool { return *a.Clear },
				Handler: func() error {
					return action2.Clear(action2.ClearOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "test", Description: "runs tests"},
				Active: func() bool { return *a.Test },
				Handler: func() error {
					return action2.Test(action2.TestOptions{
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
					return action2.Welcome(action2.WelcomeOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "help", Description: "shows the help menu"},
				Active: func() bool { return *a.Help },
				Handler: func() error {
					return action2.Help(action2.HelpOptions{})
				},
			},
			{
				Choice: search.Choice{Id: "version", Description: "shows binary version"},
				Active: func() bool { return *a.Version },
				Handler: func() error {
					return action2.Version(action2.VersionOptions{
						Efs: a.Efs,
					})
				},
			},
		},
	}, nil
}
