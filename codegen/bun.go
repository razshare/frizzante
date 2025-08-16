package codegen

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/platform"
	"github.com/razshare/frizzante/files"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
)

func Bun(c *cli.Cli, clr bool, base string) error {
	dst := filepath.Join(base, ".gen", "bun")
	plat, err := platform.Find(c, clr)
	if err != nil {
		return err
	}

	var url string

	if plat == platform.DarwinArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-aarch64.zip"
	} else if plat == platform.DarwinAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-x64.zip"
	} else if plat == platform.LinuxArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-aarch64.zip"
	} else if plat == platform.LinuxAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-x64.zip"
	} else if plat == platform.WindowsArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-windows-x64-baseline.zip"
	} else if plat == platform.WindowsAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-windows-x64-baseline.zip"
	}

	install, err := Download(url)
	if err != nil {
		return err
	}

	installed, err := install(dst)
	if err != nil {
		return err
	}

	if !installed {
		return nil
	}

	if plat == platform.DarwinArm64 {
		err = files.Move(filepath.Join(dst, "bun-darwin-aarch64", "bun"), filepath.Join(dst, "bun"))
	} else if plat == platform.DarwinAmd64 {
		err = files.Move(filepath.Join(dst, "bun-darwin-x64", "bun"), filepath.Join(dst, "bun"))
	} else if plat == platform.LinuxArm64 {
		err = files.Move(filepath.Join(dst, "bun-linux-aarch64", "bun"), filepath.Join(dst, "bun"))
	} else if plat == platform.LinuxAmd64 {
		err = files.Move(filepath.Join(dst, "bun-linux-x64", "bun"), filepath.Join(dst, "bun"))
	} else if plat == platform.WindowsArm64 {
		err = files.Move(filepath.Join(dst, "bun-windows-x64-baseline", "bun"), filepath.Join(dst, "bun.exe"))
	} else if plat == platform.WindowsAmd64 {
		err = files.Move(filepath.Join(dst, "bun-windows-x64-baseline", "bun"), filepath.Join(dst, "bun.exe"))
	}

	if err != nil {
		return err
	}

	if plat == platform.DarwinArm64 {
		err = os.Remove(filepath.Join(dst, "bun-darwin-aarch64"))
	} else if plat == platform.DarwinAmd64 {
		err = os.Remove(filepath.Join(dst, "bun-darwin-x64"))
	} else if plat == platform.LinuxArm64 {
		err = os.Remove(filepath.Join(dst, "bun-linux-aarch64"))
	} else if plat == platform.LinuxAmd64 {
		err = os.Remove(filepath.Join(dst, "bun-linux-x64"))
	} else if plat == platform.WindowsArm64 {
		err = os.Remove(filepath.Join(dst, "bun-windows-x64-baseline"))
	} else if plat == platform.WindowsAmd64 {
		err = os.Remove(filepath.Join(dst, "bun-windows-x64-baseline"))
	}

	if err != nil {
		return err
	}

	if plat != platform.WindowsArm64 && plat != platform.WindowsAmd64 {
		return unix.Chmod(".gen/bun/bun", 0755)
	}

	return nil
}
