package on

import "github.com/razshare/frizzante/cli"

func Configure(c *cli.Cli, base string) error {
	err := Generate(c, base, "bun,air")
	if err != nil {
		return err
	}
	return Install(c, base)
}
