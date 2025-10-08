package user

import (
	"runtime"
	"strings"

	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/platform"
	"github.com/razshare/frizzante/tui/messages"
)

func Platform(app *apps.App) (plat platform.Platform, err error) {
	var platStr string
	if *app.Platform != "" {
		platStr = strings.TrimSpace(*app.Platform)
	} else {
		platStr = runtime.GOOS + "/" + runtime.GOARCH
	}

	if strings.ToLower(platStr) == "linux/amd64" {
		plat = platform.LinuxAmd64
		return
	}

	if strings.ToLower(platStr) == "linux/arm64" {
		plat = platform.LinuxArm64
		return
	}

	if strings.ToLower(platStr) == "darwin/arm64" {
		plat = platform.DarwinArm64
		return
	}

	if strings.ToLower(platStr) == "darwin/amd64" {
		plat = platform.DarwinAmd64
		return
	}

	if strings.ToLower(platStr) == "windows/arm64" {
		plat = platform.WindowsArm64
		return
	}

	if strings.ToLower(platStr) == "windows/amd64" {
		plat = platform.WindowsAmd64
		return
	}

	messages.Infof("unknown platform %s, falling back to linux/amd64", platStr)

	plat = platform.LinuxAmd64

	return
}
