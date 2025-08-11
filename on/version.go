package on

import (
	"embed"
	"github.com/razshare/frizzante/tui/messages"
	"strings"
)

func Version(efs embed.FS) {
	var version string

	versionData, versionError := efs.ReadFile("version")
	if versionError != nil {
		messages.Fatal(versionError)
	}

	version = string(versionData)

	lines := strings.Split(version, "\n")

	if len(lines) == 0 {
		return
	}

	println(lines[0])
}
