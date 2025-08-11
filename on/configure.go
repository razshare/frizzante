package on

import (
	"embed"
)

func Configure(efs embed.FS) {
	Generate(efs, "bun,air")
	Install()
}
