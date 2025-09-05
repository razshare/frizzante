package generate

import (
	"path/filepath"
	"syscall"

	"github.com/razshare/frizzante/platform"
)

func Air(options AirOptions) (err error) {
	var url string

	if options.Platform == platform.DarwinArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_arm64"
	} else if options.Platform == platform.DarwinAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_amd64"
	} else if options.Platform == platform.LinuxArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_arm64"
	} else if options.Platform == platform.LinuxAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64"
	} else if options.Platform == platform.WindowsArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_arm64.exe"
	} else if options.Platform == platform.WindowsAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_amd64.exe"
	}

	var install Install
	if install, _, err = Download(DownloadOptions{Url: url, Auto: options.Auto}); err != nil {
		return
	}

	if _, err = install(filepath.Dir(options.Air)); err != nil {
		return
	}

	if options.Platform != platform.WindowsArm64 && options.Platform != platform.WindowsAmd64 && filepath.Separator != '\\' {
		err = syscall.Chmod(options.Air, 0755)
	}

	return
}
