package generations

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/razshare/frizzante/v2/cli/caches"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/platforms"
)

func Bun(options BunOptions) (err error) {
	platform := platforms.Detect()
	var url string
	if platform == platforms.DarwinArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.3.13/bun-darwin-aarch64.zip"
	} else if platform == platforms.DarwinAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.3.13/bun-darwin-x64.zip"
	} else if platform == platforms.LinuxArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.3.13/bun-linux-aarch64.zip"
	} else if platform == platforms.LinuxAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.3.13/bun-linux-x64.zip"
	} else if platform == platforms.WindowsArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.3.13/bun-windows-x64-baseline.zip"
	} else if platform == platforms.WindowsAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.3.13/bun-windows-x64-baseline.zip"
	}
	var fileName string
	if fileName, err = caches.DownloadFile(caches.DownloadFileOptions{Url: url}); err != nil {
		return
	}
	if err = caches.Install(caches.InstallOptions{
		FromFileName:    fileName,
		ToDirectoryName: filepath.Dir(options.Bun),
	}); err != nil {
		return
	}
	if platform == platforms.DarwinArm64 {
		err = files.Move(filepath.Join(filepath.Dir(options.Bun), "bun-darwin-aarch64", "bun"), options.Bun)
	} else if platform == platforms.DarwinAmd64 {
		err = files.Move(filepath.Join(filepath.Dir(options.Bun), "bun-darwin-x64", "bun"), options.Bun)
	} else if platform == platforms.LinuxArm64 {
		err = files.Move(filepath.Join(filepath.Dir(options.Bun), "bun-linux-aarch64", "bun"), options.Bun)
	} else if platform == platforms.LinuxAmd64 {
		err = files.Move(filepath.Join(filepath.Dir(options.Bun), "bun-linux-x64", "bun"), options.Bun)
	} else if platform == platforms.WindowsArm64 {
		err = files.Move(filepath.Join(filepath.Dir(options.Bun), "bun-windows-x64-baseline", "bun.exe"), options.Bun)
	} else if platform == platforms.WindowsAmd64 {
		err = files.Move(filepath.Join(filepath.Dir(options.Bun), "bun-windows-x64-baseline", "bun.exe"), options.Bun)
	}
	if err != nil {
		return
	}
	if platform == platforms.DarwinArm64 {
		err = os.Remove(filepath.Join(filepath.Dir(options.Bun), "bun-darwin-aarch64"))
	} else if platform == platforms.DarwinAmd64 {
		err = os.Remove(filepath.Join(filepath.Dir(options.Bun), "bun-darwin-x64"))
	} else if platform == platforms.LinuxArm64 {
		err = os.Remove(filepath.Join(filepath.Dir(options.Bun), "bun-linux-aarch64"))
	} else if platform == platforms.LinuxAmd64 {
		err = os.Remove(filepath.Join(filepath.Dir(options.Bun), "bun-linux-x64"))
	} else if platform == platforms.WindowsArm64 {
		err = os.Remove(filepath.Join(filepath.Dir(options.Bun), "bun-windows-x64-baseline"))
	} else if platform == platforms.WindowsAmd64 {
		err = os.Remove(filepath.Join(filepath.Dir(options.Bun), "bun-windows-x64-baseline"))
	}
	if err != nil {
		return
	}
	if platform != platforms.WindowsArm64 && platform != platforms.WindowsAmd64 && filepath.Separator != '\\' {
		err = syscall.Chmod(options.Bun, 0755)
	}
	return
}
