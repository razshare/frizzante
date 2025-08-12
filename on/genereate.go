package on

import (
	"embed"
	"github.com/razshare/frizzante/codegen"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/multiselect"
	"strings"
)

func Generate(efs embed.FS, n string) {
	if n == ":pick" {
		feats := multiselect.Send(
			[]string{
				`Core
					Generates router and view swapping tools.
				`,
				`Forms
					Generates a <Form> component that allows management of pending requests and errors.
				`,
				`Links
					Generates a <Link> component that allows management of pending requests and errors.
				`,
				`Air
					Generates Air binaries, a ☁️ Live reload tool for Go apps.
				`,
				`Bun
					Generates Bun binaries, a fast JavaScript all-in-one toolkit.
				`,
				`Session
					Generates functions for managing user session state.
				`,
				`Database
					Generates a full database setup and defaults for querying it using SQLC.
				`,
				`Queries
					Generates Go code from your ./lib/database/queries.sql file using SQLC.
				`,
			},
			"What to generate",
		)

		for _, feat := range feats {
			gen, exists := codegen.Functions[strings.ToLower(feat)]
			if !exists {
				messages.Fatalf("feature `%s` not found", feat)
			}
			gen(efs)
		}
		return
	}

	for _, feat := range strings.Split(n, ",") {
		gen, exists := codegen.Functions[strings.ToLower(feat)]
		if !exists {
			messages.Fatalf("feature `%s` not found", feat)
		}
		gen(efs)
	}
	return
}
