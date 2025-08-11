package platforms

import (
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"strings"
)

func Find() PlatformType {
	var platform string

	if *flags.Platform != "" {
		platform = *flags.Platform
	} else {
		platform = singleselect.Send(
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

		*flags.Platform = platform
	}

	if strings.ToLower(platform) == "linux/amd64" {
		return PlatformTypeLinuxAmd64
	}

	if strings.ToLower(platform) == "linux/arm64" {
		return PlatformTypeLinuxArm64
	}

	if strings.ToLower(platform) == "darwin/arm64" {
		return PlatformTypeDarwinArm64
	}

	if strings.ToLower(platform) == "darwin/amd64" {
		return PlatformTypeDarwinAmd64
	}

	if strings.ToLower(platform) == "windows/arm64" {
		return PlatformTypeWindowsArm64
	}

	if strings.ToLower(platform) == "windows/amd64" {
		return PlatformTypeWindowsAmd64
	}

	messages.Fatalf("unknown platform `%s`", platform)
	return PlatformTypeLinuxAmd64 // Noop, cli.Fatalf will crash intentionally.
}
