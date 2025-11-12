package actions

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/configs"
	"github.com/razshare/frizzante/tui/inputs"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_many"
	"github.com/razshare/frizzante/tui/select_one"
)

func Generate(options GenerateOptions) (err error) {
	pick := func(gen string) (err error) {
		fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
		fmt.Println(configs.Styles.Menu.Render(fmt.Sprintf("running ▷ generate (generates resources) ▷ %s", gen)))

		if gen == "air" {
			return generate.Air(generate.AirOptions{
				Air:      options.Air,
				Auto:     options.Auto,
				Platform: options.Platform,
			})
		} else if gen == "air_config" {
			return generate.AirConfig(generate.AirConfigOptions{
				Efs:  options.Efs,
				Tags: options.Tags,
			})
		} else if gen == "bun" {
			return generate.Bun(generate.BunOptions{
				Bun:      options.Bun,
				Auto:     options.Auto,
				Platform: options.Platform,
			})
		} else if gen == "database" {
			return generate.Database(generate.DatabaseOptions{
				Generate: gen,
				Auto:     options.Auto,
				Go:       options.Go,
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
				Efs:      options.Efs,
			})
		} else if gen == "queries" {
			return generate.Queries(generate.QueriesOptions{
				Auto:     options.Auto,
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
			})
		} else if gen == "schema" {
			var databaseString string
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
				Auto:     options.Auto,
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
				SqlcYaml: options.SqlcYaml,
				Database: database,
			})
		} else if gen == "core" {
			return generate.Core(generate.CoreOptions{
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		} else if gen == "forms" {
			return generate.Forms(generate.FormsOptions{
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		} else if gen == "links" {
			return generate.Links(generate.LinksOptions{
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		} else if gen == "icons" {
			return generate.Icons(generate.IconsOptions{
				Bun:  options.Bun,
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		} else if gen == "types" {
			return generate.TypeDefinitions(generate.TypeDefinitionsOptions{
				Go: options.Go,
			})
		} else if gen == "security" {
			return generate.Security(generate.SecurityOptions{
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		} else if gen == "migration" {
			return generate.Migration(generate.MigrationOptions{
				Auto:     options.Auto,
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
				SqlcYaml: options.SqlcYaml,
			})
		} else if gen == "sqlc" {
			return generate.Sqlc(generate.SqlcOptions{
				Auto:     options.Auto,
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
			})
		}

		return errors.New("unknown generation")
	}

	if options.Generation == "" {
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

	for _, item := range strings.Split(options.Generation, ",") {
		if err = pick(strings.ToLower(item)); err != nil {
			return
		}
	}

	return
}
