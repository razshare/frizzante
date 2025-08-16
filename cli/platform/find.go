package platform

import (
	"fmt"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/text"
	"strings"
)

func Find(c *cli.Cli, clr bool) (Type, error) {
	var plat string
	var err error

	if *c.Flags.Platform != "" {
		plat = *c.Flags.Platform
	} else {
		plat, err = singleselect.Send(
			[]string{
				"Linux/amd64",
				"Linux/arm64",
				"Darwin/amd64",
				"Darwin/arm64",
				"Windows/amd64",
				"Windows/arm64",
			},
			"platform",
		)

		*c.Flags.Platform = plat
	}

	if err != nil {
		return 0, err
	}

	if clr {
		text.Clrscr()
	}

	if strings.ToLower(plat) == "linux/amd64" {
		return LinuxAmd64, nil
	}

	if strings.ToLower(plat) == "linux/arm64" {
		return LinuxArm64, nil
	}

	if strings.ToLower(plat) == "darwin/arm64" {
		return DarwinArm64, nil
	}

	if strings.ToLower(plat) == "darwin/amd64" {
		return DarwinAmd64, nil
	}

	if strings.ToLower(plat) == "windows/arm64" {
		return WindowsArm64, nil
	}

	if strings.ToLower(plat) == "windows/amd64" {
		return WindowsAmd64, nil
	}

	return 0, fmt.Errorf("unknown platform %s", plat)
}
