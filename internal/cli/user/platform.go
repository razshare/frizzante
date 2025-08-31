package user

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/internal/cli/app"
	platform2 "github.com/razshare/frizzante/internal/platform"
	messages2 "github.com/razshare/frizzante/internal/tui/messages"
)

var PlatformMutex sync.Mutex

func Platform(a *app.App) (plat platform2.Platform, err error) {
	var home string
	if home, err = FrizzanteHome(); err != nil {
		return 0, err
	}

	name := filepath.Join(home, "platform.txt")

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
				messages2.Error(err)
				return
			}
		}

		if err = os.WriteFile(name, []byte(platStr), os.ModePerm); err != nil {
			messages2.Error(err)
		}
	}

	if strings.ToLower(platStr) == "linux/amd64" {
		save()
		return platform2.LinuxAmd64, nil
	}

	if strings.ToLower(platStr) == "linux/arm64" {
		save()
		return platform2.LinuxArm64, nil
	}

	if strings.ToLower(platStr) == "darwin/arm64" {
		save()
		return platform2.DarwinArm64, nil
	}

	if strings.ToLower(platStr) == "darwin/amd64" {
		save()
		return platform2.DarwinAmd64, nil
	}

	if strings.ToLower(platStr) == "windows/arm64" {
		save()
		return platform2.WindowsArm64, nil
	}

	if strings.ToLower(platStr) == "windows/amd64" {
		save()
		return platform2.WindowsAmd64, nil
	}

	messages2.Infof("unknown platform %s, falling back to linux/amd64", platStr)

	return platform2.LinuxAmd64, nil
}
