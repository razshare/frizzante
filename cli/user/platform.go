package user

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/razshare/frizzante/cli/app"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/platform"
	"github.com/razshare/frizzante/tui/messages"
)

var PlatformMutex sync.Mutex

func Platform(a *app.App) (plat platform.Platform, err error) {
	var cache string
	if cache, err = FrizzanteCache(); err != nil {
		return 0, err
	}

	name := filepath.Join(cache, "platform.txt")

	var platStr string
	if files.IsFile(name) {
		var data []byte
		if data, err = os.ReadFile(name); err != nil {
			return 0, err
		}

		platStr = strings.TrimSpace(string(data))
	} else {
		platStr = strings.TrimSpace(*a.Platform)
	}

	if platStr == "" {
		platStr = runtime.GOOS + "/" + runtime.GOARCH
	}

	save := func() {
		if dir := filepath.Dir(name); !files.IsDirectory(dir) {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				messages.Error(err)
				return
			}
		}

		if err = os.WriteFile(name, []byte(platStr), os.ModePerm); err != nil {
			messages.Error(err)
		}
	}

	if strings.ToLower(platStr) == "linux/amd64" {
		save()
		return platform.LinuxAmd64, nil
	}

	if strings.ToLower(platStr) == "linux/arm64" {
		save()
		return platform.LinuxArm64, nil
	}

	if strings.ToLower(platStr) == "darwin/arm64" {
		save()
		return platform.DarwinArm64, nil
	}

	if strings.ToLower(platStr) == "darwin/amd64" {
		save()
		return platform.DarwinAmd64, nil
	}

	if strings.ToLower(platStr) == "windows/arm64" {
		save()
		return platform.WindowsArm64, nil
	}

	if strings.ToLower(platStr) == "windows/amd64" {
		save()
		return platform.WindowsAmd64, nil
	}

	messages.Infof("unknown platform %s, falling back to linux/amd64", platStr)

	return platform.LinuxAmd64, nil
}
