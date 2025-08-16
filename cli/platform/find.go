package platform

import (
	"fmt"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/tui/singleselect"
	"strings"
)

func Find(c *cli.Cli) (Type, error) {
	var p string
	var err error

	if *c.Flags.Platform != "" {
		p = *c.Flags.Platform
	} else {
		p, err = singleselect.Send(
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

		*c.Flags.Platform = p
	}

	if err != nil {
		return 0, err
	}

	if strings.ToLower(p) == "linux/amd64" {
		return LinuxAmd64, nil
	}

	if strings.ToLower(p) == "linux/arm64" {
		return LinuxArm64, nil
	}

	if strings.ToLower(p) == "darwin/arm64" {
		return DarwinArm64, nil
	}

	if strings.ToLower(p) == "darwin/amd64" {
		return DarwinAmd64, nil
	}

	if strings.ToLower(p) == "windows/arm64" {
		return WindowsArm64, nil
	}

	if strings.ToLower(p) == "windows/amd64" {
		return WindowsAmd64, nil
	}

	return 0, fmt.Errorf("unknown platform %s", p)
}
