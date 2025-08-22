package codegen

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/platform"
	"os"
	"path/filepath"
	"syscall"
)

func Bun(opts BunOptions) (err error) {
	var url string

	if opts.Platform == platform.DarwinArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-aarch64.zip"
	} else if opts.Platform == platform.DarwinAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-x64.zip"
	} else if opts.Platform == platform.LinuxArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-aarch64.zip"
	} else if opts.Platform == platform.LinuxAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-x64.zip"
	} else if opts.Platform == platform.WindowsArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-windows-x64-baseline.zip"
	} else if opts.Platform == platform.WindowsAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-windows-x64-baseline.zip"
	}

	var install Install
	if install, _, err = Download(DownloadOptions{Url: url, Auto: opts.Auto}); err != nil {
		return
	}

	var installed bool
	if installed, err = install(filepath.Dir(opts.Bun)); err != nil {
		return
	}

	if !installed {
		return
	}

	if opts.Platform == platform.DarwinArm64 {
		err = files.Move(filepath.Join(filepath.Dir(opts.Bun), "bun-darwin-aarch64", "bun"), opts.Bun)
	} else if opts.Platform == platform.DarwinAmd64 {
		err = files.Move(filepath.Join(filepath.Dir(opts.Bun), "bun-darwin-x64", "bun"), opts.Bun)
	} else if opts.Platform == platform.LinuxArm64 {
		err = files.Move(filepath.Join(filepath.Dir(opts.Bun), "bun-linux-aarch64", "bun"), opts.Bun)
	} else if opts.Platform == platform.LinuxAmd64 {
		err = files.Move(filepath.Join(filepath.Dir(opts.Bun), "bun-linux-x64", "bun"), opts.Bun)
	} else if opts.Platform == platform.WindowsArm64 {
		err = files.Move(filepath.Join(filepath.Dir(opts.Bun), "bun-windows-x64-baseline", "bun.exe"), opts.Bun)
	} else if opts.Platform == platform.WindowsAmd64 {
		err = files.Move(filepath.Join(filepath.Dir(opts.Bun), "bun-windows-x64-baseline", "bun.exe"), opts.Bun)
	}

	if err != nil {
		return
	}

	if opts.Platform == platform.DarwinArm64 {
		err = os.Remove(filepath.Join(filepath.Dir(opts.Bun), "bun-darwin-aarch64"))
	} else if opts.Platform == platform.DarwinAmd64 {
		err = os.Remove(filepath.Join(filepath.Dir(opts.Bun), "bun-darwin-x64"))
	} else if opts.Platform == platform.LinuxArm64 {
		err = os.Remove(filepath.Join(filepath.Dir(opts.Bun), "bun-linux-aarch64"))
	} else if opts.Platform == platform.LinuxAmd64 {
		err = os.Remove(filepath.Join(filepath.Dir(opts.Bun), "bun-linux-x64"))
	} else if opts.Platform == platform.WindowsArm64 {
		err = os.Remove(filepath.Join(filepath.Dir(opts.Bun), "bun-windows-x64-baseline"))
	} else if opts.Platform == platform.WindowsAmd64 {
		err = os.Remove(filepath.Join(filepath.Dir(opts.Bun), "bun-windows-x64-baseline"))
	}

	if err != nil {
		return
	}

	if opts.Platform != platform.WindowsArm64 && opts.Platform != platform.WindowsAmd64 && filepath.Separator != '\\' {
		err = syscall.Chmod(opts.Bun, 0755)
	}

	return
}
