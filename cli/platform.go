package cli

import (
	"github.com/pterm/pterm"
	"strings"
)

func Platform() PlatformType {
	var platform string

	if *FlagPlatform != "" {
		platform = *FlagPlatform
	} else {
		var platformError error
		platform, platformError = pterm.
			DefaultInteractiveSelect.
			WithOptions([]string{
				"Linux/amd64",
				"Linux/arm64",
				"Darwin/amd64",
				"Darwin/arm64",
				"Windows/amd64",
				"Windows/arm64",
			}).
			Show("Pick a platform")

		if platformError != nil {
			Fatal(platformError)
		}
		*FlagPlatform = platform
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

	Fatalf("unknown platform `%s`", platform)
	return PlatformTypeLinuxAmd64 // Noop, cli.Fatalf will crash intentionally.
}
