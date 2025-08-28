package codegen

import (
	"github.com/razshare/frizzante/platform"
	"path/filepath"
)

func Sqlc(options SqlcOptions) (err error) {
	var url string

	if options.Platform == platform.DarwinArm64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_arm64.zip"
	} else if options.Platform == platform.DarwinAmd64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_amd64.zip"
	} else if options.Platform == platform.LinuxArm64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_arm64.zip"
	} else if options.Platform == platform.LinuxAmd64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_amd64.zip"
	} else if options.Platform == platform.WindowsArm64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
	} else if options.Platform == platform.WindowsAmd64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
	}

	var install Install
	if install, _, err = Download(DownloadOptions{Url: url, Auto: options.Auto}); err != nil {
		return
	}

	_, err = install(filepath.Dir(options.Sqlc))

	return
}
