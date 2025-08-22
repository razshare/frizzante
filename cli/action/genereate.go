package action

import (
	"errors"
	"github.com/razshare/frizzante/cli/codegen"
	"github.com/razshare/frizzante/tui/multiselect"
	"github.com/razshare/frizzante/tui/search"
	"path/filepath"
	"strings"
)

func Generate(opts GenerateOptions) error {
	err := codegen.Init(opts.Efs)
	if err != nil {
		return err
	}

	pick := func(gen string) error {
		if gen == "air" {
			return codegen.Air(codegen.AirOptions{
				Air:      opts.Air,
				Auto:     opts.Auto,
				Platform: opts.Platform,
			})
		} else if gen == "bun" {
			return codegen.Bun(codegen.BunOptions{
				Bun:      opts.Bun,
				Auto:     opts.Auto,
				Platform: opts.Platform,
			})
		} else if gen == "session" {
			return codegen.Session(codegen.SessionOptions{
				Auto: opts.Auto,
				Lib:  filepath.Join("lib", "session"),
				Efs:  opts.Efs,
			})
		} else if gen == "database" {
			return codegen.Database(codegen.DatabaseOptions{
				Generate: gen,
				Auto:     opts.Auto,
				Go:       opts.Go,
				Sqlc:     opts.Sqlc,
				Platform: opts.Platform,
				Efs:      opts.Efs,
				Lib:      filepath.Join("lib", "database"),
			})
		} else if gen == "queries" {
			return codegen.Queries(codegen.QueriesOptions{
				Auto:     opts.Auto,
				Sqlc:     opts.Sqlc,
				Platform: opts.Platform,
				Lib:      filepath.Join("lib", "database"),
			})
		} else if gen == "core" {
			return codegen.Core(codegen.CoreOptions{
				App:  opts.App,
				Auto: opts.Auto,
				Efs:  opts.Efs,
				Lib:  filepath.Join(opts.App, "frizzante", "core"),
			})
		} else if gen == "forms" {
			return codegen.Forms(codegen.FormsOptions{
				App:  opts.App,
				Auto: opts.Auto,
				Efs:  opts.Efs,
				Lib:  filepath.Join(opts.App, "frizzante", "forms"),
			})
		} else if gen == "links" {
			return codegen.Links(codegen.LinksOptions{
				App:  opts.App,
				Auto: opts.Auto,
				Efs:  opts.Efs,
				Lib:  filepath.Join(opts.App, "frizzante", "forms"),
			})
		}

		return errors.New("unknown generation")
	}

	if opts.Selected == "" {
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
			err = pick(strings.ToLower(item))
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, item := range strings.Split(opts.Selected, ",") {
		err = pick(strings.ToLower(item))
		if err != nil {
			return err
		}
	}

	return nil
}
