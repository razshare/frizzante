package codegen

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/platform"
	"os"
	"path/filepath"
	"syscall"
)

func Bun(options BunOptions) (err error) {
	var url string

	if options.Platform == platform.DarwinArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-aarch64.zip"
	} else if options.Platform == platform.DarwinAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-x64.zip"
	} else if options.Platform == platform.LinuxArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-aarch64.zip"
	} else if options.Platform == platform.LinuxAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-x64.zip"
	} else if options.Platform == platform.WindowsArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-windows-x64-baseline.zip"
	} else if options.Platform == platform.WindowsAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-windows-x64-baseline.zip"
	}

	var install Install
	if install, _, err = Download(DownloadOptions{Url: url, Auto: options.Auto}); err != nil {
		return
	}

	var installed bool
	if installed, err = install(filepath.Dir(options.Bun)); err != nil {
		return
	}

	if !installed {
		return
	}

	if options.Platform == platform.DarwinArm64 {
		err = files.Move(filepath.Join(filepath.Dir(options.Bun), "bun-darwin-aarch64", "bun"), options.Bun)
	} else if options.Platform == platform.DarwinAmd64 {
		err = files.Move(filepath.Join(filepath.Dir(options.Bun), "bun-darwin-x64", "bun"), options.Bun)
	} else if options.Platform == platform.LinuxArm64 {
		err = files.Move(filepath.Join(filepath.Dir(options.Bun), "bun-linux-aarch64", "bun"), options.Bun)
	} else if options.Platform == platform.LinuxAmd64 {
		err = files.Move(filepath.Join(filepath.Dir(options.Bun), "bun-linux-x64", "bun"), options.Bun)
	} else if options.Platform == platform.WindowsArm64 {
		err = files.Move(filepath.Join(filepath.Dir(options.Bun), "bun-windows-x64-baseline", "bun.exe"), options.Bun)
	} else if options.Platform == platform.WindowsAmd64 {
		err = files.Move(filepath.Join(filepath.Dir(options.Bun), "bun-windows-x64-baseline", "bun.exe"), options.Bun)
	}

	if err != nil {
		return
	}

	if options.Platform == platform.DarwinArm64 {
		err = os.Remove(filepath.Join(filepath.Dir(options.Bun), "bun-darwin-aarch64"))
	} else if options.Platform == platform.DarwinAmd64 {
		err = os.Remove(filepath.Join(filepath.Dir(options.Bun), "bun-darwin-x64"))
	} else if options.Platform == platform.LinuxArm64 {
		err = os.Remove(filepath.Join(filepath.Dir(options.Bun), "bun-linux-aarch64"))
	} else if options.Platform == platform.LinuxAmd64 {
		err = os.Remove(filepath.Join(filepath.Dir(options.Bun), "bun-linux-x64"))
	} else if options.Platform == platform.WindowsArm64 {
		err = os.Remove(filepath.Join(filepath.Dir(options.Bun), "bun-windows-x64-baseline"))
	} else if options.Platform == platform.WindowsAmd64 {
		err = os.Remove(filepath.Join(filepath.Dir(options.Bun), "bun-windows-x64-baseline"))
	}

	if err != nil {
		return
	}

	if options.Platform != platform.WindowsArm64 && options.Platform != platform.WindowsAmd64 && filepath.Separator != '\\' {
		err = syscall.Chmod(options.Bun, 0755)
	}

	return
}
