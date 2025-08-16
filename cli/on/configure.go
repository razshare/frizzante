package on

import "github.com/razshare/frizzante/cli"

func Configure(c *cli.Cli, clr bool, base string) error {
	err := Generate(c, clr, base, "bun,air")
	if err != nil {
		return err
	}
	return Install(c, base)
}
