package on

import "embed"

func Configure(efs embed.FS, base string) error {
	err := Generate(efs, base, "bun,air")
	if err != nil {
		return err
	}
	return Install(base)
}
