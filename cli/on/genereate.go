package on

import (
	"fmt"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/codegen"
	"github.com/razshare/frizzante/tui/multiselect"
	"github.com/razshare/frizzante/tui/text"
	"strings"
)

func Generate(c *cli.Cli, clr bool, base string, gen string) error {
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
			generate, exists := codegen.Functions[strings.ToLower(item)]
			if !exists {
				return fmt.Errorf("unknown option %s", item)
			}

			err = generate(c, false, base)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, item := range strings.Split(gen, ",") {
		generate, exists := codegen.Functions[strings.ToLower(item)]
		if !exists {
			return fmt.Errorf("unknown option %s", item)
		}
		err := generate(c, false, base)
		if err != nil {
			return err
		}
	}
	return nil
}
