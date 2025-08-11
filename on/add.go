package on

import (
	"embed"
	"github.com/razshare/frizzante/gen"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/multiselect"
	"strings"
)

func Add(efs embed.FS, n string) {
	if n == ":pick" {
		feats := multiselect.Send(
			[]string{
				"Core",
				"Forms",
				"Links",
				"Air",
				"Bun",
			},
			"Pick a feature to add",
		)

		for _, feat := range feats {
			gen, exists := codegen.Generators[feat]
			if !exists {
				messages.Fatalf("feature `%s` not found", feat)
			}
			gen(efs)
		}
		return
	}

	for _, feat := range strings.Split(n, ",") {
		gen, exists := codegen.Generators[feat]
		if !exists {
			messages.Fatalf("feature `%s` not found", feat)
		}
		gen(efs)
	}
	return
}
