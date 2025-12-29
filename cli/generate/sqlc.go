package generate

import (
	"path/filepath"

	"github.com/razshare/frizzante/cli/caches"
	"github.com/razshare/frizzante/platforms"
)

func Sqlc(options SqlcOptions) (err error) {
	platform := platforms.Detect()

	var url string
	if platform == platforms.DarwinArm64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_arm64.zip"
	} else if platform == platforms.DarwinAmd64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_amd64.zip"
	} else if platform == platforms.LinuxArm64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_arm64.zip"
	} else if platform == platforms.LinuxAmd64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_amd64.zip"
	} else if platform == platforms.WindowsArm64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
	} else if platform == platforms.WindowsAmd64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
	}

	var fileName string
	if fileName, err = caches.DownloadFile(caches.DownloadFileOptions{
		Url: url,
	}); err != nil {
		return
	}

	err = caches.Install(caches.InstallOptions{
		FromFileName:    fileName,
		ToDirectoryName: filepath.Dir(options.Sqlc),
	})

	return
}
