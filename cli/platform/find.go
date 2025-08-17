package platform

import (
	"fmt"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/text"
	"strings"
)

func Find(clr bool) (Type, error) {
	var plat string
	var err error

	if *cli.Platform != "" {
		plat = *cli.Platform
	} else {
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

		*cli.Platform = plat
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
