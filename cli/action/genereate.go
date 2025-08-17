package action

import (
	"errors"
	"github.com/razshare/frizzante/cli/codegen"
	"github.com/razshare/frizzante/tui/multiselect"
	"path/filepath"
	"strings"
)

func Generate(o GenerateOptions) error {
	pick := func(gen string) error {
		if gen == "air" {
			return codegen.Air(codegen.AirOptions{
				Air:      o.Air,
				Auto:     o.Auto,
				Platform: o.Platform,
			})
		} else if gen == "bun" {
			return codegen.Bun(codegen.BunOptions{
				Bun:      o.Bun,
				Auto:     o.Auto,
				Platform: o.Platform,
			})
		} else if gen == "session" {
			return codegen.Session(codegen.SessionOptions{
				Auto: o.Auto,
				Lib:  filepath.Join("lib", "session"),
			})
		} else if gen == "database" {
			return codegen.Database(codegen.DatabaseOptions{
				Generate: gen,
				Auto:     o.Auto,
				Go:       o.Go,
				Sqlc:     o.Sqlc,
				Platform: o.Platform,
				Lib:      filepath.Join("lib", "database"),
			})
		} else if gen == "queries" {
			return codegen.Queries(codegen.QueriesOptions{
				Auto:     o.Auto,
				Sqlc:     o.Sqlc,
				Platform: o.Platform,
				Lib:      filepath.Join("lib", "database"),
			})
		} else if gen == "core" {
			return codegen.Core(codegen.CoreOptions{
				App:  o.App,
				Auto: o.Auto,
				Lib:  filepath.Join(o.App, "frizzante", "core"),
			})
		} else if gen == "forms" {
			return codegen.Forms(codegen.FormsOptions{
				App:  o.App,
				Auto: o.Auto,
				Lib:  filepath.Join(o.App, "frizzante", "forms"),
			})
		} else if gen == "links" {
			return codegen.Links(codegen.LinksOptions{
				App:  o.App,
				Auto: o.Auto,
				Lib:  filepath.Join(o.App, "frizzante", "forms"),
			})
		}

		return errors.New("unknown generation")
	}

	if o.Selected == ":pick" {
		items, err := multiselect.Send(
			[]string{
				`core
					router and view swapping tools.
				`,
				`forms
					form component that provides status details
				`,
				`links
					hyperlink component that provides status details
				`,
				`air
					live reload tool for go programs
				`,
				`bun
					fast js toolkit
				`,
				`session
					functions for managing user session state
				`,
				`database
					full database setup
				`,
				`queries
					sql code to go code using sqlc
				`,
			},
			"generate",
		)

		if err != nil {
			return err
		}

		for _, item := range items {
			err = pick(strings.ToLower(item))
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, item := range strings.Split(o.Selected, ",") {
		err := pick(strings.ToLower(item))
		if err != nil {
			return err
		}
	}
	return nil
}
