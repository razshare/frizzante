package cli

import "embed"

func OnConfigure(efs embed.FS) {
	OnAddFeature(efs, "bun,air")
	OnInstall()
}
