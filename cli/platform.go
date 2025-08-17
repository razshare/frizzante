package cli

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/platform"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

func Platform(c *Cli) (platform.Platform, error) {
	var plat string

	home, err := os.UserHomeDir()
	if err != nil {
		return 0, err
	}

	ntxt := filepath.Join(home, ".frizzante", "platform.txt")

	if files.IsFile(ntxt) {
		var d []byte
		d, err = os.ReadFile(ntxt)
		if err != nil {
			return 0, err
		}

		plat = strings.TrimSpace(string(d))
	} else {
		plat = strings.TrimSpace(*c.Platform)
	}

	if plat == "" {
		detected := runtime.GOOS + "/" + runtime.GOARCH
		if *c.Yes {
			plat = detected
		} else {
			list := []string{
				"linux/amd64",
				"linux/arm64",
				"darwin/amd64",
				"darwin/arm64",
				"windows/amd64",
				"windows/arm64",
			}

			if i := slices.Index(list, detected); i >= 0 {
				if reduced := append(list[:i], list[i+1:]...); reduced != nil {
					list = append([]string{detected}, reduced...)
				}
			}

			plat, err = singleselect.Send(list, "platform (detected "+detected+")")
		}
	}

	if err != nil {
		return 0, err
	}

	save := func() {
		dir := filepath.Dir(ntxt)
		if !files.IsDirectory(dir) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				messages.Error(err)
				return
			}
		}

		err = os.WriteFile(ntxt, []byte(plat), os.ModePerm)
		if err != nil {
			messages.Error(err)
		}
	}

	if strings.ToLower(plat) == "linux/amd64" {
		save()
		return platform.PlatformLinuxAmd64, nil
	}

	if strings.ToLower(plat) == "linux/arm64" {
		save()
		return platform.PlatformLinuxArm64, nil
	}

	if strings.ToLower(plat) == "darwin/arm64" {
		save()
		return platform.PlatformDarwinArm64, nil
	}

	if strings.ToLower(plat) == "darwin/amd64" {
		save()
		return platform.PlatformDarwinAmd64, nil
	}

	if strings.ToLower(plat) == "windows/arm64" {
		save()
		return platform.PlatformWindowsArm64, nil
	}

	if strings.ToLower(plat) == "windows/amd64" {
		save()
		return platform.PlatformWindowsAmd64, nil
	}

	messages.Infof("unknown platform %s, falling back to linux/amd64", plat)

	return platform.PlatformLinuxAmd64, nil
}
