package action

import (
	"errors"
	"github.com/razshare/frizzante/cli/codegen"
	"github.com/razshare/frizzante/tui/multiselect"
	"github.com/razshare/frizzante/tui/search"
	"path/filepath"
	"strings"
)

func Generate(options GenerateOptions) (err error) {
	if err = codegen.Init(options.Efs); err != nil {
		return
	}

	pick := func(gen string) error {
		if gen == "air" {
			return codegen.Air(codegen.AirOptions{
				Air:      options.Air,
				Auto:     options.Auto,
				Platform: options.Platform,
			})
		} else if gen == "bun" {
			return codegen.Bun(codegen.BunOptions{
				Bun:      options.Bun,
				Auto:     options.Auto,
				Platform: options.Platform,
			})
		} else if gen == "session" {
			return codegen.Session(codegen.SessionOptions{
				Auto: options.Auto,
				Lib:  filepath.Join("lib", "session"),
				Efs:  options.Efs,
			})
		} else if gen == "database" {
			return codegen.Database(codegen.DatabaseOptions{
				Generate: gen,
				Auto:     options.Auto,
				Go:       options.Go,
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
				Efs:      options.Efs,
				Lib:      filepath.Join("lib", "database"),
			})
		} else if gen == "queries" {
			return codegen.Queries(codegen.QueriesOptions{
				Auto:     options.Auto,
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
				Lib:      filepath.Join("lib", "database"),
			})
		} else if gen == "core" {
			return codegen.Core(codegen.CoreOptions{
				App:  options.App,
				Auto: options.Auto,
				Efs:  options.Efs,
				Lib:  filepath.Join(options.App, "frizzante", "core"),
			})
		} else if gen == "forms" {
			return codegen.Forms(codegen.FormsOptions{
				App:  options.App,
				Auto: options.Auto,
				Efs:  options.Efs,
				Lib:  filepath.Join(options.App, "frizzante", "forms"),
			})
		} else if gen == "links" {
			return codegen.Links(codegen.LinksOptions{
				App:  options.App,
				Auto: options.Auto,
				Efs:  options.Efs,
				Lib:  filepath.Join(options.App, "frizzante", "forms"),
			})
		}

		return errors.New("unknown generation")
	}

	if options.Selected == "" {
		var items []string
		items, err = multiselect.Send(
			[]search.Choice{
				{Id: "core", Description: "router and view swapping tools."},
				{Id: "forms", Description: "form component that provides status details"},
				{Id: "links", Description: "hyperlink component that provides status details"},
				{Id: "air", Description: "live reload tool for go programs"},
				{Id: "bun", Description: "fast js toolkit"},
				{Id: "session", Description: "functions for managing user session state"},
				{Id: "database", Description: "full database setup"},
				{Id: "queries", Description: "sql code to go code using sqlc"},
			},
			"generate",
		)

		if err != nil {
			return err
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
