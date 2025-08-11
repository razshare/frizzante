package on

import (
	"embed"
)

func Configure(efs embed.FS) {
	Add(efs, "bun,air")
	Install()
}
