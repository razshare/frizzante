package cli

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/platform"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"os"
	"path/filepath"
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
		plat, err = singleselect.Send(
			[]string{
				"linux/amd64",
				"linux/arm64",
				"darwin/amd64",
				"darwin/arm64",
				"windows/amd64",
				"windows/arm64",
			},
			"platform",
		)
	}

	if err != nil {
		return 0, err
	}

	save := func() {
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

	messages.Errorf("unknown platform %s", plat)

	return Platform(c)
}
