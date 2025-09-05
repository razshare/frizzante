package action

import (
	"errors"
	"strings"

	generate2 "github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/tui/multiselect"
	"github.com/razshare/frizzante/tui/search"
)

func Generate(options GenerateOptions) (err error) {
	pick := func(gen string) error {
		if gen == "air" {
			return generate2.Air(generate2.AirOptions{
				Air:      options.Air,
				Auto:     options.Auto,
				Platform: options.Platform,
			})
		} else if gen == "bun" {
			return generate2.Bun(generate2.BunOptions{
				Bun:      options.Bun,
				Auto:     options.Auto,
				Platform: options.Platform,
			})
		} else if gen == "session" {
			return generate2.Session(generate2.SessionOptions{
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		} else if gen == "database" {
			return generate2.Database(generate2.DatabaseOptions{
				Generate: gen,
				Auto:     options.Auto,
				Go:       options.Go,
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
				Efs:      options.Efs,
			})
		} else if gen == "queries" {
			return generate2.Queries(generate2.QueriesOptions{
				Auto:     options.Auto,
				Sqlc:     options.Sqlc,
				Platform: options.Platform,
			})
		} else if gen == "core" {
			return generate2.Core(generate2.CoreOptions{
				App:  options.App,
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		} else if gen == "forms" {
			return generate2.Forms(generate2.FormsOptions{
				App:  options.App,
				Auto: options.Auto,
				Efs:  options.Efs,
			})
		} else if gen == "links" {
			return generate2.Links(generate2.LinksOptions{
				App:  options.App,
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
				{Id: "core", Description: "server, routing and view swapping tools."},
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
