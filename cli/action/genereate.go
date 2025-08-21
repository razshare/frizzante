package action

import (
	"errors"
	"github.com/razshare/frizzante/cli/codegen"
	"github.com/razshare/frizzante/tui/multiselect"
	"github.com/razshare/frizzante/tui/search"
	"path/filepath"
	"strings"
)

func Generate(o GenerateOptions) error {
	err := codegen.Init(o.Efs)
	if err != nil {
		return err
	}

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
				Efs:  o.Efs,
			})
		} else if gen == "database" {
			return codegen.Database(codegen.DatabaseOptions{
				Generate: gen,
				Auto:     o.Auto,
				Go:       o.Go,
				Sqlc:     o.Sqlc,
				Platform: o.Platform,
				Efs:      o.Efs,
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
				Efs:  o.Efs,
				Lib:  filepath.Join(o.App, "frizzante", "core"),
			})
		} else if gen == "forms" {
			return codegen.Forms(codegen.FormsOptions{
				App:  o.App,
				Auto: o.Auto,
				Efs:  o.Efs,
				Lib:  filepath.Join(o.App, "frizzante", "forms"),
			})
		} else if gen == "links" {
			return codegen.Links(codegen.LinksOptions{
				App:  o.App,
				Auto: o.Auto,
				Efs:  o.Efs,
				Lib:  filepath.Join(o.App, "frizzante", "forms"),
			})
		}

		return errors.New("unknown generation")
	}

	if o.Selected == "" {
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

	for _, item := range strings.Split(o.Selected, ",") {
		err = pick(strings.ToLower(item))
		if err != nil {
			return err
		}
	}

	return nil
}
