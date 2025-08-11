package platform

import (
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"strings"
)

func Find() Type {
	var p string

	if *flags.Platform != "" {
		p = *flags.Platform
	} else {
		p = singleselect.Send(
			[]string{
				"Linux/amd64",
				"Linux/arm64",
				"Darwin/amd64",
				"Darwin/arm64",
				"Windows/amd64",
				"Windows/arm64",
			},
			"Pick a platform",
		)

		*flags.Platform = p
	}

	if strings.ToLower(p) == "linux/amd64" {
		return LinuxAmd64
	}

	if strings.ToLower(p) == "linux/arm64" {
		return LinuxArm64
	}

	if strings.ToLower(p) == "darwin/arm64" {
		return DarwinArm64
	}

	if strings.ToLower(p) == "darwin/amd64" {
		return DarwinAmd64
	}

	if strings.ToLower(p) == "windows/arm64" {
		return WindowsArm64
	}

	if strings.ToLower(p) == "windows/amd64" {
		return WindowsAmd64
	}

	messages.Fatalf("unknown platform `%s`", p)
	return LinuxAmd64 // Noop, messages.Fatalf will crash.
}
