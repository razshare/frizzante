package codegen

import (
	"github.com/razshare/frizzante/platform"
	"path/filepath"
	"syscall"
)

func Air(opts AirOptions) (err error) {
	var url string

	if opts.Platform == platform.DarwinArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_arm64"
	} else if opts.Platform == platform.DarwinAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_amd64"
	} else if opts.Platform == platform.LinuxArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_arm64"
	} else if opts.Platform == platform.LinuxAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64"
	} else if opts.Platform == platform.WindowsArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_arm64.exe"
	} else if opts.Platform == platform.WindowsAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_amd64.exe"
	}

	var install Install
	if install, _, err = Download(DownloadOptions{Url: url, Auto: opts.Auto}); err != nil {
		return
	}

	if _, err = install(filepath.Dir(opts.Air)); err != nil {
		return
	}

	if opts.Platform != platform.WindowsArm64 && opts.Platform != platform.WindowsAmd64 && filepath.Separator != '\\' {
		err = syscall.Chmod(opts.Air, 0755)
	}

	return
}
