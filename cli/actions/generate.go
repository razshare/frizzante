package actions

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/configs"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/inputs"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_many"
	"github.com/razshare/frizzante/tui/select_one"
)

func Generate(options GenerateOptions) (err error) {
	generation := options.Generation
	if generation == ":pick" {
		generation = ""
	}

	pick := func(gen string) (err error) {
		fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
		fmt.Println(configs.Styles.Menu.Render(fmt.Sprintf("running ▷ generate (generates resources) ▷ %s", gen)))

		if gen == "air" {
			directoryName := filepath.Dir(options.Air)
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
				Air:      options.Air,
				Platform: options.Platform,
			})
		} else if gen == "air_config" {
			if err = generate.AirConfig(generate.AirConfigOptions{
				Efs:  options.Efs,
				Tags: options.Tags,
			}); err != nil {
				return err
			}
		} else if gen == "bun" {
			directoryName := filepath.Dir(options.Bun)
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
				Bun:      options.Bun,
				Platform: options.Platform,
			})
		} else if gen == "sqlc" {
			directoryName := filepath.Dir(options.Sqlc)
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
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
			})
		} else if gen == "database" {
			databaseType := options.DatabaseType
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

			return generate.Database(generate.DatabaseOptions{
				Generate: gen,
				Go:       options.Go,
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
				Efs:      options.Efs,
				Type:     databaseType,
			})
		} else if gen == "migration" {
			yamlFileName := options.SqlcYaml
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

			return generate.Migration(generate.MigrationOptions{
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
				SqlcYaml: yamlFileName,
			})
		} else if gen == "queries" {
			yamlFileName := options.SqlcYaml
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

			return generate.Queries(generate.QueriesOptions{
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
				SqlcYaml: yamlFileName,
			})
		} else if gen == "schema" {
			yamlFileName := options.SqlcYaml
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

			databaseString := options.Database
			if options.Database == "" {
				var names []string
				if names, err = files.FindWithSuffix("lib", ".sqlite"); err != nil {
					return
				}

				choices := make([]search.Choice, len(names))
				for index, name := range names {
					choices[index] = search.Choice{Id: name}
				}

				choices = append(choices, search.Choice{Id: "other", Description: "use a different file"})

				if databaseString, err = select_one.Sendf(choices, "where's your sqlite database located?"); err != nil {
					return
				}

				if databaseString == "other" {
					if databaseString, err = inputs.Send("where's the file located?"); err != nil {
						return
					}
				}
			} else {
				databaseString = options.Database
			}

			var database *sql.DB
			if database, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", databaseString)); err != nil {
				return
			}

			return generate.Schema(generate.SchemaOptions{
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
				SqlcYaml: options.SqlcYaml,
				Database: database,
			})
		} else if gen == "core" {
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

			return generate.Core(generate.CoreOptions{Efs: options.Efs})
		} else if gen == "forms" {
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
			return generate.Forms(generate.FormsOptions{Efs: options.Efs})
		} else if gen == "links" {
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
			return generate.Links(generate.LinksOptions{
				Efs: options.Efs,
			})
		} else if gen == "icons" {
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
			return generate.Icons(generate.IconsOptions{
				Bun: options.Bun,
				Efs: options.Efs,
			})
		} else if gen == "security" {
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
			return generate.Security(generate.SecurityOptions{Efs: options.Efs})
		} else if gen == "types" {
			return generate.TypeDefinitions(generate.TypeDefinitionsOptions{
				Go: options.Go,
			})
		}

		return errors.New("unknown generation")
	}

	if generation == "" {
		var items []string
		items, err = select_many.Send(
			[]search.Choice{
				{Id: "core", Description: "core features"},
				{Id: "forms", Description: "form component that provides status details"},
				{Id: "links", Description: "hyperlink component that provides status details"},
				{Id: "icons", Description: "icon component that renders using svg"},
				{Id: "air", Description: "live reload tool for go programs"},
				{Id: "air_config", Description: "air configuration file .air.toml"},
				{Id: "bun", Description: "fast js toolkit"},
				{Id: "database", Description: "full database setup"},
				{Id: "sqlc", Description: "sql compiler"},
				{Id: "queries", Description: "sql code to go code using sqlc"},
				{Id: "schema", Description: "database schema"},
				{Id: "security", Description: "security and cryptographic functions"},
				{Id: "types", Description: "typescript type definitions using .d.ts files"},
				{Id: "migration", Description: "migration file using current date"},
			},
			"generate",
		)

		if err != nil {
			return
		}

		for _, item := range items {
			if err = pick(strings.ToLower(item)); err != nil {
				return
			}
		}

		if len(items) != 0 {
			err = Generate(options)
		}

		return
	}

	for _, item := range strings.Split(generation, ",") {
		if err = pick(strings.ToLower(item)); err != nil {
			return
		}
	}

	return
}
