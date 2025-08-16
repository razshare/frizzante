package on

import (
	"github.com/razshare/frizzante/cli"
	"strings"
)

func Version(c *cli.Cli) error {
	var v string

	d, err := c.Efs.ReadFile("version")
	if err != nil {
		return err
	}

	v = string(d)

	ls := strings.Split(v, "\n")

	if len(ls) == 0 {
		return nil
	}

	println(ls[0])

	return nil
}
