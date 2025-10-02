package user

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	app_ "github.com/razshare/frizzante/cli/app"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/platform"
	"github.com/razshare/frizzante/tui/messages"
)

var PlatformMutex sync.Mutex

func Platform(a *app_.App) (plat platform.Platform, err error) {
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

	save := func() error {
		if dir := filepath.Dir(name); !files.IsDirectory(dir) {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {

				return err
			}
		}

		if files.IsFile(name) {
			if err = os.Remove(name); err != nil {
				return err
			}
		}

		if err = os.WriteFile(name, []byte(platStr), os.ModePerm); err != nil {
			return err
		}

		return nil
	}

	if strings.ToLower(platStr) == "linux/amd64" {
		err = save()
		plat = platform.LinuxAmd64
		return
	}

	if strings.ToLower(platStr) == "linux/arm64" {
		err = save()
		plat = platform.LinuxArm64
		return
	}

	if strings.ToLower(platStr) == "darwin/arm64" {
		err = save()
		plat = platform.DarwinArm64
		return
	}

	if strings.ToLower(platStr) == "darwin/amd64" {
		err = save()
		plat = platform.DarwinAmd64
		return
	}

	if strings.ToLower(platStr) == "windows/arm64" {
		err = save()
		plat = platform.WindowsArm64
		return
	}

	if strings.ToLower(platStr) == "windows/amd64" {
		err = save()
		plat = platform.WindowsAmd64
		return
	}

	messages.Infof("unknown platform %s, falling back to linux/amd64", platStr)

	plat = platform.LinuxAmd64

	return
}
