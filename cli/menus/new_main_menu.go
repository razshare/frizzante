package menus

import (
	"database/sql"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
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

func NewMainMenu(options NewMainMenuOptions) (*Menu, error) {
	app := options.App
	efs := app.Efs
	tags := *app.Tags
	sqlcYaml := *app.SqlcYaml
	persistent := options.Persistent
	platform := platforms.Detect()

	cache, err := paths.Cache()
	if err != nil {
		return nil, err
	}

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

	var logo string
	if logo, err = apps.Logo(&app); err != nil {
		messages.Warning(err)
		err = nil
	}

	return &Menu{
		Title:      "main",
		Logo:       logo,
		Persistent: persistent,
		Items: []Item{
			{
				Ids:    []string{"configure"},
				Choice: search.Choice{Id: "configure", Description: "generates bun and air binaries"},
				Handler: func(_ string) (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ configure"))

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
						Efs:      efs,
						Platform: platform,
					})

					return
				},
			},
			{
				Ids:    []string{"create", "new"},
				Choice: search.Choice{Id: "create new project", Description: "creates a new project"},
				Handler: func(value string) (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ create project"))

					var name string
					if name = value; name == "" {
						name, err = inputs.Send("give the project a name")
						if err != nil {
							return
						}
					}

					err = actions.CreateProject(actions.CreateProjectOptions{
						Name: name,
						Go:   go_,
						Efs:  efs,
						Air:  air,
						Bun:  bun,
					})
					return err
				},
			},
			{
				Ids:    []string{"install"},
				Choice: search.Choice{Id: "install", Description: "installs go and js packages"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ install"))
					return actions.Install(actions.InstallOptions{
						Go:  go_,
						Bun: bun,
					})
				},
			},
			{
				Ids:    []string{"update"},
				Choice: search.Choice{Id: "update", Description: "updates go and js packages"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ update"))
					return actions.Update(actions.UpdateOptions{
						Go:  go_,
						Bun: bun,
					})
				},
			},
			{
				Ids:    []string{"add"},
				Choice: search.Choice{Id: "add", Description: "adds packages"},
				Handler: func(value string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ add"))
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
							Query: value,
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
				Ids:    []string{"dev"},
				Choice: search.Choice{Id: "dev", Description: "runs air and vite in parallel"},
				Handler: func(_ string) (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ dev"))
					err = actions.Dev(actions.DevOptions{
						Go:   go_,
						Air:  air,
						Bun:  bun,
						Tags: tags,
						Efs:  efs,
					})
					return
				},
			},
			{
				Ids:    []string{"build"},
				Choice: search.Choice{Id: "build", Description: "builds project"},
				Handler: func(_ string) (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ build"))
					err = actions.Build(actions.BuildOptions{
						Go:   go_,
						Bun:  bun,
						Tags: tags,
					})
					return
				},
			},
			{
				Ids:    []string{"assembly-explorer"},
				Choice: search.Choice{Id: "assembly explorer", Description: "explores application assembly output"},
				Handler: func(_ string) (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ assembly explorer"))
					err = actions.AssemblyExplorer(actions.AssemblyExplorerOptions{
						Go:   go_,
						Bun:  bun,
						Tags: tags,
					})
					return
				},
			},
			{
				Ids:    []string{"generate"},
				Choice: search.Choice{Id: "generate", Description: "generates code and resources"},
				Handler: func(value string) (err error) {
					var menu *Menu
					if menu, err = NewGenerateMenu(NewGenerateMenuOptions{
						App:        app,
						Persistent: persistent,
					}); err != nil {
						return
					}

					err = ParseQueryAndActivate(menu, value)

					return
				},
			},
			{
				Ids:    []string{"migrate"},
				Choice: search.Choice{Id: "migrate", Description: "migrates database schema"},
				Handler: func(value string) (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ migrate"))

					if sqlcYaml == "" {
						var items []string
						if items, err = files.ReadDirectory("lib"); err != nil {
							return
						}

						names := make([]string, 0)
						for _, item := range items {
							if strings.HasSuffix(item, string(filepath.Separator)+"sqlc.yaml") {
								names = append(names, item)
							}
						}

						choices := make([]search.Choice, len(names))
						for index, name := range names {
							choices[index] = search.Choice{Id: name}
						}

						choices = append(choices, search.Choice{Id: "other", Description: "use a different file"})
						sqlcYaml, err = select_one.Sendf(choices, "where is your sqlc.yaml file located?")
					}

					var offset string
					var target string

					migrateRange := strings.SplitN(value, ",", 2)

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
						Sqlc:     sqlc,
						SqlcYaml: sqlcYaml,
						Offset:   offset,
						Target:   target,
						Database: database,
						Platform: platform,
					})
					return
				},
			},
			{
				Ids:    []string{"package"},
				Choice: search.Choice{Id: "package", Description: "builds app"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ package"))
					return actions.Package(actions.PackageOptions{
						Bun: bun,
					})
				},
			},
			{
				Ids:    []string{"package-watch"},
				Choice: search.Choice{Id: "package watch", Description: "builds app on change"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ package watch"))
					return actions.PackageWatch(actions.PackageWatchOptions{
						Bun: bun,
					})
				},
			},
			{
				Ids:    []string{"check"},
				Choice: search.Choice{Id: "check", Description: "checks for code errors"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ check"))
					return actions.Check(actions.CheckOptions{
						Bun: bun,
					})
				},
			},
			{
				Ids:    []string{"format"},
				Choice: search.Choice{Id: "format", Description: "format code"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ format"))
					return actions.Format(actions.FormatOptions{
						Go:  go_,
						Bun: bun,
					})
				},
			},
			{
				Ids:    []string{"touch"},
				Choice: search.Choice{Id: "touch", Description: "adds placeholders in app/dist"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ touch"))
					return actions.Touch(actions.TouchOptions{})
				},
			},
			{
				Ids:    []string{"clean"},
				Choice: search.Choice{Id: "clean project", Description: "deletes .gen, .vite, app/{dist,node_modules}"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ clean project"))
					return actions.CleanProject(actions.CleanProjectOptions{
						Go: go_,
					})
				},
			},
			{
				Ids:    []string{"reset"},
				Choice: search.Choice{Id: "reset", Description: "deletes " + cache},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ reset"))
					return actions.Reset(actions.ResetOptions{})
				},
			},
			{
				Ids:    []string{"clear"},
				Choice: search.Choice{Id: "clear", Description: "clears screen"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ clear"))
					return actions.Clear(actions.ClearOptions{})
				},
			},
			{
				Ids:    []string{"lock-packages"},
				Choice: search.Choice{Id: "lock packages", Description: "locks packages to the current exact version"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ lock packages"))
					return actions.LockPackages(actions.LockPackagesOptions{})
				},
			},
			{
				Ids:    []string{"snapshot"},
				Choice: search.Choice{Id: "snapshot", Description: "snapshots the server state and generates static web assets"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ snapshot"))
					return actions.LockPackages(actions.LockPackagesOptions{})
				},
			},
			{
				Ids:    []string{"test"},
				Choice: search.Choice{Id: "test", Description: "runs tests"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ test"))
					return actions.Test(actions.TestOptions{
						Go:  go_,
						Bun: bun,
					})
				},
			},
			{
				Hidden: true,
				Ids:    []string{"welcome"},
				Choice: search.Choice{Id: "welcome", Description: "shows a welcome message"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ welcome"))
					return actions.Welcome(actions.WelcomeOptions{})
				},
			},
			{
				Hidden: true,
				Ids:    []string{"help"},
				Choice: search.Choice{Id: "help", Description: "shows the help menu"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ help"))
					return actions.Help(actions.HelpOptions{})
				},
			},
			{
				Ids:    []string{"version"},
				Choice: search.Choice{Id: "version", Description: "shows binary version"},
				Handler: func(_ string) error {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("running ▷ version"))
					return actions.Version(actions.VersionOptions{
						Efs: efs,
					})
				},
			},
		},
	}, nil
}
