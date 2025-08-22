package user

import (
	"github.com/razshare/frizzante/cli/app"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/platform"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Platform(a *app.App) (platform.Platform, error) {
	var plat string

	home, err := FrizzanteHome()
	if err != nil {
		return 0, err
	}

	txt := filepath.Join(home, "platform.txt")

	if files.IsFile(txt) {
		var d []byte
		d, err = os.ReadFile(txt)
		if err != nil {
			return 0, err
		}

		plat = strings.TrimSpace(string(d))
	} else {
		plat = strings.TrimSpace(*a.Platform)
	}

	if plat == "" {
		plat = runtime.GOOS + "/" + runtime.GOARCH
	}

	save := func() {
		dir := filepath.Dir(txt)
		if !files.IsDirectory(dir) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				messages.Error(err)
				return
			}
		}

		err = os.WriteFile(txt, []byte(plat), os.ModePerm)
		if err != nil {
			messages.Error(err)
		}
	}

	if strings.ToLower(plat) == "linux/amd64" {
		save()
		return platform.LinuxAmd64, nil
	}

	if strings.ToLower(plat) == "linux/arm64" {
		save()
		return platform.LinuxArm64, nil
	}

	if strings.ToLower(plat) == "darwin/arm64" {
		save()
		return platform.DarwinArm64, nil
	}

	if strings.ToLower(plat) == "darwin/amd64" {
		save()
		return platform.DarwinAmd64, nil
	}

	if strings.ToLower(plat) == "windows/arm64" {
		save()
		return platform.WindowsArm64, nil
	}

	if strings.ToLower(plat) == "windows/amd64" {
		save()
		return platform.WindowsAmd64, nil
	}

	messages.Infof("unknown platform %s, falling back to linux/amd64", plat)

	return platform.LinuxAmd64, nil
}
