package on

import (
	"embed"
	"strings"
)

func Version(efs embed.FS) error {
	var v string

	d, err := efs.ReadFile("version")
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
