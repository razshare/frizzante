package on

import (
	"errors"
	"github.com/razshare/frizzante/cli/codegen"
	"github.com/razshare/frizzante/tui/multiselect"
	"github.com/razshare/frizzante/tui/text"
	"strings"
)

func Generate(app string, gen string, clr bool, yes bool, gobin string, sqlcbin string) error {
	pick := func(gen string) error {
		if gen == "air" {
			return codegen.Air(clr)
		} else if gen == "bun" {
			return codegen.Bun(clr)
		} else if gen == "session" {
			return codegen.Session(clr, yes)
		} else if gen == "database" {
			return codegen.Database(gen, clr, yes, gobin, sqlcbin)
		} else if gen == "queries" {
			return codegen.Queries(clr, sqlcbin)
		} else if gen == "core" {
			return codegen.Core(app, yes)
		} else if gen == "forms" {
			return codegen.Forms(app, yes)
		} else if gen == "links" {
			return codegen.Links(app, yes)
		}

		return errors.New("unknown generation")
	}

	if gen == ":pick" {
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
			"what to generate",
		)

		if err != nil {
			return err
		}

		if clr {
			text.Clrscr()
		}

		for _, item := range items {
			err = pick(strings.ToLower(item))
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, item := range strings.Split(gen, ",") {
		err := pick(strings.ToLower(item))
		if err != nil {
			return err
		}
	}
	return nil
}
