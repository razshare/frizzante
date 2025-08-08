package cli

import (
	"embed"
	"strings"
)

func OnVersion(efs embed.FS) {
	var version string

	versionData, versionError := efs.ReadFile("version")
	if versionError != nil {
		Fatal(versionError)
	}

	version = string(versionData)

	lines := strings.Split(version, "\n")

	if len(lines) == 0 {
		return
	}

	println(lines[0])
}
