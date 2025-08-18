package codegen

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/platform"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
)

func Bun(o BunOptions) error {
	var url string

	if o.Platform == platform.PlatformDarwinArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-aarch64.zip"
	} else if o.Platform == platform.PlatformDarwinAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-x64.zip"
	} else if o.Platform == platform.PlatformLinuxArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-aarch64.zip"
	} else if o.Platform == platform.PlatformLinuxAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-x64.zip"
	} else if o.Platform == platform.PlatformWindowsArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-windows-x64-baseline.zip"
	} else if o.Platform == platform.PlatformWindowsAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-windows-x64-baseline.zip"
	}

	install, err := Download(DownloadOptions{
		Url:  url,
		Auto: o.Auto,
	})

	if err != nil {
		return err
	}

	installed, err := install(filepath.Dir(o.Bun))
	if err != nil {
		return err
	}

	if !installed {
		return nil
	}

	if o.Platform == platform.PlatformDarwinArm64 {
		err = files.Move(filepath.Join(filepath.Dir(o.Bun), "bun-darwin-aarch64", "bun"), o.Bun)
	} else if o.Platform == platform.PlatformDarwinAmd64 {
		err = files.Move(filepath.Join(filepath.Dir(o.Bun), "bun-darwin-x64", "bun"), o.Bun)
	} else if o.Platform == platform.PlatformLinuxArm64 {
		err = files.Move(filepath.Join(filepath.Dir(o.Bun), "bun-linux-aarch64", "bun"), o.Bun)
	} else if o.Platform == platform.PlatformLinuxAmd64 {
		err = files.Move(filepath.Join(filepath.Dir(o.Bun), "bun-linux-x64", "bun"), o.Bun)
	} else if o.Platform == platform.PlatformWindowsArm64 {
		err = files.Move(filepath.Join(filepath.Dir(o.Bun), "bun-windows-x64-baseline", "bun"), o.Bun)
	} else if o.Platform == platform.PlatformWindowsAmd64 {
		err = files.Move(filepath.Join(filepath.Dir(o.Bun), "bun-windows-x64-baseline", "bun"), o.Bun)
	}

	if err != nil {
		return err
	}

	if o.Platform == platform.PlatformDarwinArm64 {
		err = os.Remove(filepath.Join(filepath.Dir(o.Bun), "bun-darwin-aarch64"))
	} else if o.Platform == platform.PlatformDarwinAmd64 {
		err = os.Remove(filepath.Join(filepath.Dir(o.Bun), "bun-darwin-x64"))
	} else if o.Platform == platform.PlatformLinuxArm64 {
		err = os.Remove(filepath.Join(filepath.Dir(o.Bun), "bun-linux-aarch64"))
	} else if o.Platform == platform.PlatformLinuxAmd64 {
		err = os.Remove(filepath.Join(filepath.Dir(o.Bun), "bun-linux-x64"))
	} else if o.Platform == platform.PlatformWindowsArm64 {
		err = os.Remove(filepath.Join(filepath.Dir(o.Bun), "bun-windows-x64-baseline"))
	} else if o.Platform == platform.PlatformWindowsAmd64 {
		err = os.Remove(filepath.Join(filepath.Dir(o.Bun), "bun-windows-x64-baseline"))
	}

	if err != nil {
		return err
	}

	if o.Platform != platform.PlatformWindowsArm64 && o.Platform != platform.PlatformWindowsAmd64 && filepath.Separator != '\\' {
		err = unix.Chmod(o.Bun, 0755)
	}

	return nil
}
