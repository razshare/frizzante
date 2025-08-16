package codegen

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/platform"
	"golang.org/x/sys/unix"
	"path/filepath"
)

func Air(c *cli.Cli, clr bool, base string) error {
	dst := filepath.Join(base, ".gen", "air")
	plat, err := platform.Find(c, clr)
	if err != nil {
		return err
	}

	var url string

	if plat == platform.DarwinArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_arm64"
	} else if plat == platform.DarwinAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_amd64"
	} else if plat == platform.LinuxArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_arm64"
	} else if plat == platform.LinuxAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64"
	} else if plat == platform.WindowsArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_arm64.exe"
	} else if plat == platform.WindowsAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_amd64.exe"
	}

	install, err := Download(url)
	if err != nil {
		return err
	}

	_, err = install(dst)
	if err != nil {
		return err
	}

	if plat != platform.WindowsArm64 && plat != platform.WindowsAmd64 {
		err = unix.Chmod(".gen/air/air", 0755)
	}

	if err != nil {
		return err
	}

	return nil
}
