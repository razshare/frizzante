package on

import (
	"embed"
	"github.com/razshare/frizzante/tui/messages"
	"strings"
)

func Version(efs embed.FS) {
	var v string

	d, err := efs.ReadFile("version")
	if err != nil {
		messages.Fatal(err)
	}

	v = string(d)

	ls := strings.Split(v, "\n")

	if len(ls) == 0 {
		return
	}

	println(ls[0])
}
