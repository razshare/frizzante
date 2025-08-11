package on

import (
	"embed"
	"github.com/razshare/frizzante/codegen/generators"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/multiselect"
	"strings"
)

func Generate(efs embed.FS, n string) {
	if n == ":pick" {
		feats := multiselect.Send(
			[]string{
				"Core",
				"Forms",
				"Links",
				"Air",
				"Bun",
				"Session",
				"Database",
			},
			"Pick a feature to add",
		)

		for _, feat := range feats {
			gen, exists := generators.Functions[strings.ToLower(feat)]
			if !exists {
				messages.Fatalf("feature `%s` not found", feat)
			}
			gen(efs)
		}
		return
	}

	for _, feat := range strings.Split(n, ",") {
		gen, exists := generators.Functions[strings.ToLower(feat)]
		if !exists {
			messages.Fatalf("feature `%s` not found", feat)
		}
		gen(efs)
	}
	return
}
