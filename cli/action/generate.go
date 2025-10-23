package action

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/razshare/frizzante/cli/generate"
	tags_ "github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/input"
	"github.com/razshare/frizzante/tui/multiselect"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
)

func Generate(options GenerateOptions) (err error) {
	pick := func(gen string) (err error) {
		if gen == "air" {
			return generate.Air(generate.AirOptions{
				Air:      options.Air,
				Auto:     options.Auto,
				Platform: options.Platform,
			})
		} else if gen == "air_config" {
			tags := make([]string, 0)

			if !options.Active {
				if tags, err = tags_.Select([]search.Choice{
					{Id: "trace", Description: "enables tracing with stack.Trace()"},
					{Id: "types", Description: "enables type generations"},
					{Id: "no_js_runtime", Description: "disables the server-side JavaScript runtime"},
					{Id: "experimental_qjs_runtime", Description: "replaces goja with qjs"},
					{Id: "other", Description: "adds custom tags"},
				}); err != nil {
					return err
				}
			}

			tags = append(tags, options.Tags...)

			return generate.AirConfig(generate.AirConfigOptions{
				Efs:  options.Efs,
				Tags: tags,
			})
		} else if gen == "bun" {
			return generate.Bun(generate.BunOptions{
				Bun:      options.Bun,
				Auto:     options.Auto,
				Platform: options.Platform,
			})
		} else if gen == "sessions" {
			return generate.Sessions(generate.SessionsOptions{
				Auto: options.Auto,
				Efs:  options.Efs,
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
				SqlcYaml: options.SqlcYaml,
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

				if databaseString, err = singleselect.Sendf(choices, "where's your sqlite database located?"); err != nil {
					return
				}

				if databaseString == "other" {
					if databaseString, err = input.Send("where's the file located?"); err != nil {
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
				App:  options.App,
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		} else if gen == "forms" {
			return generate.Forms(generate.FormsOptions{
				App:  options.App,
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		} else if gen == "links" {
			return generate.Links(generate.LinksOptions{
				App:  options.App,
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		} else if gen == "icons" {
			return generate.Icons(generate.IconsOptions{
				App:  options.App,
				Bun:  options.Bun,
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		} else if gen == "types" {
			return generate.Types(generate.TypesOptions{
				Auto: options.Auto,
				Efs:  options.Efs,
				Go:   options.Go,
			})
		} else if gen == "security" {
			return generate.Security(generate.SecurityOptions{
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		}

		return errors.New("unknown generation")
	}

	if options.Selected == "" {
		var items []string
		items, err = multiselect.Send(
			[]search.Choice{
				{Id: "core", Description: "core features"},
				{Id: "forms", Description: "form component that provides status details"},
				{Id: "links", Description: "hyperlink component that provides status details"},
				{Id: "icons", Description: "icon component that renders using svg"},
				{Id: "air", Description: "live reload tool for go programs"},
				{Id: "air_config", Description: "air configuration file .air.toml"},
				{Id: "bun", Description: "fast js toolkit"},
				{Id: "sessions", Description: "features for managing user sessions"},
				{Id: "database", Description: "full database setup"},
				{Id: "queries", Description: "sql code to go code using sqlc"},
				{Id: "schema", Description: "update the schema"},
				{Id: "security", Description: "security and cryptographic functions"},
				{Id: "types", Description: "type definitions using .d.ts files"},
			},
			"generate",
		)

		if slices.Contains(items, "types:features") && slices.Contains(items, "types") {
			typesFeaturesIndex := slices.Index(items, "types:features")
			typesIndex := slices.Index(items, "types")

			if typesIndex < typesFeaturesIndex {
				items[typesIndex] = "types:features"
				items[typesFeaturesIndex] = "types"
			}
		}

		if err != nil {
			return
		}

		for _, item := range items {
			if err = pick(strings.ToLower(item)); err != nil {
				return
			}
		}

		return
	}

	for _, item := range strings.Split(options.Selected, ",") {
		if err = pick(strings.ToLower(item)); err != nil {
			return
		}
	}

	return
}
