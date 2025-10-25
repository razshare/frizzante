package platforms

import (
	"runtime"
	"strings"

	"github.com/razshare/frizzante/tui/messages"
)

// Detect automatically detects the current platform.
//
// If the current platform doesn't match any of the allowed platforms, Detect falls back to Linux.
func Detect() (plat Platform) {
	platform := runtime.GOOS + "/" + runtime.GOARCH

	if strings.ToLower(platform) == "linux/amd64" {
		plat = LinuxAmd64
		return
	}

	if strings.ToLower(platform) == "linux/arm64" {
		plat = LinuxArm64
		return
	}

	if strings.ToLower(platform) == "darwin/arm64" {
		plat = DarwinArm64
		return
	}

	if strings.ToLower(platform) == "darwin/amd64" {
		plat = DarwinAmd64
		return
	}

	if strings.ToLower(platform) == "windows/arm64" {
		plat = WindowsArm64
		return
	}

	if strings.ToLower(platform) == "windows/amd64" {
		plat = WindowsAmd64
		return
	}

	messages.Infof("unknown platform %s, falling back to linux/amd64", platform)

	plat = LinuxAmd64

	return
}
