package on

import (
	"github.com/razshare/frizzante/cli"
	"strings"
)

func Version() error {
	var v string

	d, err := cli.Efs.ReadFile("version")
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
