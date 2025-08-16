package on

import (
	"embed"
	"fmt"
	"github.com/razshare/frizzante/codegen"
	"github.com/razshare/frizzante/tui/multiselect"
	"strings"
)

func Generate(efs embed.FS, base string, n string) error {
	if n == ":pick" {
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

		for _, item := range items {
			generate, exists := codegen.Functions[strings.ToLower(item)]
			if !exists {
				return fmt.Errorf("unknown option %s", item)
			}

			err = generate(efs, base)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, item := range strings.Split(n, ",") {
		generate, exists := codegen.Functions[strings.ToLower(item)]
		if !exists {
			return fmt.Errorf("unknown option %s", item)
		}
		generate(efs, base)
	}
	return nil
}
