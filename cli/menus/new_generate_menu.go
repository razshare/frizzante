package menus

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/cli/paths"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/platforms"
	"github.com/razshare/frizzante/tui/configs"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/inputs"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
)

func NewGenerateMenu(options NewGenerateMenuOptions) (*Menu, error) {
	app := options.App
	efs := app.Efs
	tags := *app.Tags
	databaseType := *app.DatabaseType
	sqlcYaml := *app.SqlcYaml
	databaseConnectionString := *app.DatabaseConnectionString
	persistent := options.Persistent
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
		Title:      "generate",
		Persistent: persistent,
		Items: []Item{
			{
				Ids:    []string{"air"},
				Choice: search.Choice{Id: "air", Description: "shows binary version"},
				Handler: func(_ string) (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("generate ▷ air"))
					directoryName := filepath.Dir(air)
					if files.IsDirectory(directoryName) {
						var yesRemove bool
						if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", directoryName); err != nil {
							return
						}

						if !yesRemove {
							err = errors.New("cannot continue generating air binaries because they already exist")
						}

						if err = os.RemoveAll(directoryName); err != nil {
							return
						}
					}

					return generate.Air(generate.AirOptions{
						Air:      air,
						Platform: platform,
					})
				},
			},
			{
				Ids:    []string{"air-config"},
				Choice: search.Choice{Id: "air config", Description: "shows binary version"},
				Handler: func(_ string) (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("generate ▷ air config"))
					if err = generate.AirConfig(generate.AirConfigOptions{
						Efs:  efs,
						Tags: tags,
					}); err != nil {
						return
					}
					return
				},
			},
			{
				Ids:    []string{"bun"},
				Choice: search.Choice{Id: "bun", Description: "shows binary version"},
				Handler: func(_ string) (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("generate ▷ air config"))

					directoryName := filepath.Dir(bun)
					if files.IsDirectory(directoryName) {
						var yesRemove bool
						if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", directoryName); err != nil {
							return
						}

						if !yesRemove {
							err = errors.New("cannot continue generating bun binaries because they already exist")
							return
						}

						if err = os.RemoveAll(directoryName); err != nil {
							return
						}
					}

					return generate.Bun(generate.BunOptions{
						Bun:      bun,
						Platform: platform,
					})
				},
			},
			{
				Ids:    []string{"sqlc"},
				Choice: search.Choice{Id: "sqlc", Description: "shows binary version"},
				Handler: func(_ string) (err error) {
					fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
					fmt.Println(configs.Styles.Menu.Render("generate ▷ air config"))

					directoryName := filepath.Dir(sqlc)
					if files.IsDirectory(directoryName) {
						var yesRemove bool
						if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", directoryName); err != nil {
							return
						}

						if !yesRemove {
							err = errors.New("cannot continue generating bun binaries because they already exist")
							return
						}

						if err = os.RemoveAll(directoryName); err != nil {
							return
						}
					}

					return generate.Sqlc(generate.SqlcOptions{
						Sqlc:     sqlc,
						Platform: platform,
					})
				},
			},
			{
				Ids:    []string{"database"},
				Choice: search.Choice{Id: "database", Description: "shows binary version"},
				Handler: func(value string) (err error) {
					if databaseType == "" {
						if databaseType, err = select_one.Send(
							[]search.Choice{{Id: "sqlite"}},
							"what type of database would you like to setup?",
						); err != nil {
							return
						}
					}

					directoryName := filepath.Join("lib", databaseType, "databases")
					var yesRemove bool
					if yesRemove, err = confirm.Sendf(true, "%s already exists. Overwrite?", directoryName); err != nil {
						return
					}

					if !yesRemove {
						err = errors.New("cannot continue generating database files because they already exist")
						return nil
					}

					if err = os.RemoveAll(directoryName); err != nil {
						return
					}

					err = generate.Database(generate.DatabaseOptions{
						Go:       go_,
						Sqlc:     sqlc,
						Platform: platform,
						Efs:      efs,
						Type:     databaseType,
					})
					return
				},
			},
			{
				Ids:    []string{"migration"},
				Choice: search.Choice{Id: "migration", Description: "shows binary version"},
				Handler: func(value string) (err error) {
					yamlFileName := sqlcYaml
					if yamlFileName == "" {
						var names []string
						if names, err = files.FindWithSuffix("lib", "sqlc.yaml"); err != nil {
							return
						}

						choices := make([]search.Choice, len(names))
						for index, name := range names {
							choices[index] = search.Choice{Id: name}
						}

						choices = append(choices, search.Choice{Id: "other", Description: "other"})

						if len(choices) == 0 {
							err = errors.New("cannot continue generating queries because no sqlc.yaml file has been provided")
							return
						}

						yamlFileName, err = select_one.Sendf(choices, "where is your sqlc.yaml file located?")
					}

					err = generate.Migration(generate.MigrationOptions{
						Sqlc:     sqlc,
						Platform: platform,
						SqlcYaml: yamlFileName,
					})

					return
				},
			},
			{
				Ids:    []string{"queries"},
				Choice: search.Choice{Id: "queries", Description: "shows binary version"},
				Handler: func(value string) (err error) {
					yamlFileName := sqlcYaml
					if yamlFileName == "" {
						var names []string
						if names, err = files.FindWithSuffix("lib", "sqlc.yaml"); err != nil {
							return
						}

						choices := make([]search.Choice, len(names))
						for index, name := range names {
							choices[index] = search.Choice{Id: name}
						}

						choices = append(choices, search.Choice{Id: "other", Description: "other"})

						if len(choices) == 0 {
							err = errors.New("cannot continue generating queries because no sqlc.yaml file has been provided")
							return
						}

						yamlFileName, err = select_one.Sendf(choices, "where is your sqlc.yaml file located?")
					}

					err = generate.Queries(generate.QueriesOptions{
						Sqlc:     sqlc,
						Platform: platform,
						SqlcYaml: yamlFileName,
					})
					return
				},
			},
			{
				Ids:    []string{"schema"},
				Choice: search.Choice{Id: "schema", Description: "shows binary version"},
				Handler: func(value string) (err error) {
					if sqlcYaml == "" {
						var names []string
						if names, err = files.FindWithSuffix("lib", "sqlc.yaml"); err != nil {
							return
						}

						choices := make([]search.Choice, len(names))
						for index, name := range names {
							choices[index] = search.Choice{Id: name}
						}

						choices = append(choices, search.Choice{Id: "other", Description: "other"})

						if len(choices) == 0 {
							err = errors.New("cannot continue generating queries because no sqlc.yaml file has been provided")
							return
						}

						sqlcYaml, err = select_one.Sendf(choices, "where is your sqlc.yaml file located?")
					}

					if databaseConnectionString == "" {
						var names []string
						if names, err = files.FindWithSuffix("lib", ".sqlite"); err != nil {
							return
						}

						choices := make([]search.Choice, len(names))
						for index, name := range names {
							choices[index] = search.Choice{Id: name}
						}

						choices = append(choices, search.Choice{Id: "other", Description: "use a different file"})

						if databaseConnectionString, err = select_one.Sendf(choices, "where's your sqlite database located?"); err != nil {
							return
						}

						if databaseConnectionString == "other" {
							if databaseConnectionString, err = inputs.Send("where's the file located?"); err != nil {
								return
							}
						}
					}

					var database *sql.DB
					if database, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", databaseConnectionString)); err != nil {
						return
					}

					err = generate.Schema(generate.SchemaOptions{
						Sqlc:     sqlc,
						Platform: platform,
						SqlcYaml: sqlcYaml,
						Database: database,
					})
					return
				},
			},
			{
				Ids:    []string{"core"},
				Choice: search.Choice{Id: "core", Description: "shows binary version"},
				Handler: func(value string) (err error) {

					libCoreDirectoryName := filepath.Join("lib", "core")
					if files.IsDirectory(libCoreDirectoryName) {
						var yesRemove bool
						if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", libCoreDirectoryName); err != nil {
							return err
						}

						if yesRemove {
							if err = os.RemoveAll(libCoreDirectoryName); err != nil {
								return
							}
						}
					}

					appLibCoreScriptsCoreDirectoryName := filepath.Join("app", "lib", "scripts", "core")
					if files.IsDirectory(appLibCoreScriptsCoreDirectoryName) {
						var yesRemove bool
						if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", appLibCoreScriptsCoreDirectoryName); err != nil {
							return err
						}

						if yesRemove {
							if err = os.RemoveAll(appLibCoreScriptsCoreDirectoryName); err != nil {
								return
							}
						}
					}

					appLibComponentsCoreDirectoryName := filepath.Join("app", "lib", "components", "core")
					if files.IsDirectory(appLibComponentsCoreDirectoryName) {
						var yesRemove bool
						if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", appLibComponentsCoreDirectoryName); err != nil {
							return err
						}

						if yesRemove {
							if err = os.RemoveAll(appLibComponentsCoreDirectoryName); err != nil {
								return
							}
						}
					}

					err = generate.Core(generate.CoreOptions{Efs: efs})
					return
				},
			},
			{
				Ids:    []string{"forms"},
				Choice: search.Choice{Id: "forms", Description: "shows binary version"},
				Handler: func(value string) (err error) {
					appLibComponentsFormsDirectoryName := filepath.Join("app", "lib", "components", "forms")
					if files.IsDirectory(appLibComponentsFormsDirectoryName) {
						var yesRemove bool
						if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", appLibComponentsFormsDirectoryName); err != nil {
							return err
						}

						if yesRemove {
							if err = os.RemoveAll(appLibComponentsFormsDirectoryName); err != nil {
								return
							}
						}
					}
					err = generate.Forms(generate.FormsOptions{Efs: efs})
					return
				},
			},
			{
				Ids:    []string{"links"},
				Choice: search.Choice{Id: "links", Description: "shows binary version"},
				Handler: func(value string) (err error) {
					appLibComponentsLinksDirectoryName := filepath.Join("app", "lib", "components", "links")
					if files.IsDirectory(appLibComponentsLinksDirectoryName) {
						var yesRemove bool
						if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", appLibComponentsLinksDirectoryName); err != nil {
							return err
						}

						if yesRemove {
							if err = os.RemoveAll(appLibComponentsLinksDirectoryName); err != nil {
								return
							}
						}
					}
					err = generate.Links(generate.LinksOptions{
						Efs: efs,
					})
					return
				},
			},
			{
				Ids:    []string{"icons"},
				Choice: search.Choice{Id: "icons", Description: "shows binary version"},
				Handler: func(value string) (err error) {
					appLibComponentsIconsDirectoryName := filepath.Join("app", "lib", "components", "icons")
					if files.IsDirectory(appLibComponentsIconsDirectoryName) {
						var yesRemove bool
						if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", appLibComponentsIconsDirectoryName); err != nil {
							return err
						}

						if yesRemove {
							if err = os.RemoveAll(appLibComponentsIconsDirectoryName); err != nil {
								return
							}
						}
					}
					err = generate.Icons(generate.IconsOptions{
						Bun: bun,
						Efs: efs,
					})
					return
				},
			},
			{
				Ids:    []string{"security"},
				Choice: search.Choice{Id: "security", Description: "shows binary version"},
				Handler: func(value string) (err error) {
					libSecurityDirectoryName := filepath.Join("lib", "security")
					if files.IsDirectory(libSecurityDirectoryName) {
						var yesRemove bool
						if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", libSecurityDirectoryName); err != nil {
							return err
						}

						if yesRemove {
							if err = os.RemoveAll(libSecurityDirectoryName); err != nil {
								return
							}
						}
					}
					err = generate.Security(generate.SecurityOptions{Efs: efs})
					return
				},
			},
			{
				Ids:    []string{"types"},
				Choice: search.Choice{Id: "types", Description: "shows binary version"},
				Handler: func(value string) (err error) {
					err = generate.TypeDefinitions(generate.TypeDefinitionsOptions{
						Go: go_,
					})
					return
				},
			},
		},
	}, nil
}
