package action

import (
	"errors"
	"strings"

	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/tui/multiselect"
	"github.com/razshare/frizzante/tui/search"
)

func Generate(options GenerateOptions) (err error) {
	pick := func(gen string) error {
		if gen == "air" {
			return generate.Air(generate.AirOptions{
				Air:      options.Air,
				Auto:     options.Auto,
				Platform: options.Platform,
			})
		} else if gen == "bun" {
			return generate.Bun(generate.BunOptions{
				Bun:      options.Bun,
				Auto:     options.Auto,
				Platform: options.Platform,
			})
		} else if gen == "session" {
			return generate.Session(generate.SessionOptions{
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
