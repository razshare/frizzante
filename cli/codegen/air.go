package codegen

import (
	"github.com/razshare/frizzante/platform"
	"golang.org/x/sys/unix"
	"path/filepath"
)

func Air(o AirOptions) error {
	var url string

	if o.Platform == platform.PlatformDarwinArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_arm64"
	} else if o.Platform == platform.PlatformDarwinAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_amd64"
	} else if o.Platform == platform.PlatformLinuxArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_arm64"
	} else if o.Platform == platform.PlatformLinuxAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64"
	} else if o.Platform == platform.PlatformWindowsArm64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_arm64.exe"
	} else if o.Platform == platform.PlatformWindowsAmd64 {
		url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_amd64.exe"
	}

	install, err := Download(DownloadOptions{
		Url:  url,
		Auto: o.Auto,
	})

	if err != nil {
		return err
	}

	_, err = install(filepath.Dir(o.Air))
	if err != nil {
		return err
	}

	if o.Platform != platform.PlatformWindowsArm64 && o.Platform != platform.PlatformWindowsAmd64 && filepath.Separator != '\\' {
		err = unix.Chmod(o.Air, 0755)
	}

	if err != nil {
		return err
	}

	return nil
}
