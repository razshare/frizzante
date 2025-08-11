package download

import (
	"embed"
	"github.com/razshare/frizzante/cli/platforms"
	"path/filepath"
)

func Air(efs embed.FS) {
	dn := filepath.Join(".gen", "air")
	p := platforms.Find()

	var url string

	if p == platforms.PlatformTypeDarwinArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_arm64"
	} else if p == platforms.PlatformTypeDarwinAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_amd64"
	} else if p == platforms.PlatformTypeLinuxArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_arm64"
	} else if p == platforms.PlatformTypeLinuxAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64"
	} else if p == platforms.PlatformTypeWindowsArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_arm64.exe"
	} else if p == platforms.PlatformTypeWindowsAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_amd64.exe"
	}

	Install("air", url, dn)

	return
}
