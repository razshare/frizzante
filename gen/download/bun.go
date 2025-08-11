package download

import (
	"embed"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/platforms"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func Bun(efs embed.FS) {
	dn := filepath.Join(".gen", "bun")
	p := platforms.Find()

	var url string

	if p == platforms.PlatformTypeDarwinArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-aarch64.zip"
	} else if p == platforms.PlatformTypeDarwinAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-x64.zip"
	} else if p == platforms.PlatformTypeLinuxArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-aarch64.zip"
	} else if p == platforms.PlatformTypeLinuxAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-x64.zip"
	} else if p == platforms.PlatformTypeWindowsArm64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-windows-x64-baseline.zip"
	} else if p == platforms.PlatformTypeWindowsAmd64 {
		url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-windows-x64-baseline.zip"
	}

	Install("bun", url, dn)

	var n string

	if p == platforms.PlatformTypeDarwinArm64 {
		n = filepath.Join(dn, "bun-darwin-aarch64", "bun")
	} else if p == platforms.PlatformTypeDarwinAmd64 {
		n = filepath.Join(dn, "bun-darwin-x64", "bun")
	} else if p == platforms.PlatformTypeLinuxArm64 {
		n = filepath.Join(dn, "bun-linux-aarch64", "bun")
	} else if p == platforms.PlatformTypeLinuxAmd64 {
		n = filepath.Join(dn, "bun-linux-x64", "bun")
	} else if p == platforms.PlatformTypeWindowsArm64 {
		n = filepath.Join(dn, "bun-windows-x64-baseline", "bun.exe")
	} else if p == platforms.PlatformTypeWindowsAmd64 {
		n = filepath.Join(dn, "bun-windows-x64-baseline", "bun.exe")
	}

	if files.IsFile(n) {
		err := os.Rename(n, path.Bun("."))
		if err != nil {
			messages.Fatal(err)
		}
	}

	idn := filepath.Dir(n)
	if files.IsDirectory(dn) {
		err := os.RemoveAll(idn)
		if err != nil {
			messages.Fatal(err)
		}
	}
}
