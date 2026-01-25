package generations

import (
	"path/filepath"
	"syscall"

	"github.com/razshare/frizzante/cli/caches"
	"github.com/razshare/frizzante/platforms"
)

func Air(options AirOptions) (err error) {
	platform := platforms.Detect()
	var url string
	if platform == platforms.DarwinArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_arm64"
	} else if platform == platforms.DarwinAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_amd64"
	} else if platform == platforms.LinuxArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_arm64"
	} else if platform == platforms.LinuxAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64"
	} else if platform == platforms.WindowsArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_arm64.exe"
	} else if platform == platforms.WindowsAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_amd64.exe"
	}
	var fileName string
	if fileName, err = caches.DownloadFile(caches.DownloadFileOptions{Url: url}); err != nil {
		return
	}
	if err = caches.Install(caches.InstallOptions{
		FromFileName:    fileName,
		ToDirectoryName: filepath.Dir(options.Air),
	}); err != nil {
		return
	}
	if platform != platforms.WindowsArm64 && platform != platforms.WindowsAmd64 && filepath.Separator != '\\' {
		err = syscall.Chmod(options.Air, 0755)
	}
	return
}
