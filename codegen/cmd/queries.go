package cmd

import (
	"embed"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/platform"
	"github.com/razshare/frizzante/codegen/download"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"path/filepath"
)

func Queries(efs embed.FS) {
	if !files.IsFile(path.Sqlc(".")) {
		dn := filepath.Join(".gen", "sqlc")
		p := platform.Find()

		var url string

		if p == platform.DarwinArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_arm64.zip"
		} else if p == platform.DarwinAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_amd64.zip"
		} else if p == platform.LinuxArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_arm64.zip"
		} else if p == platform.LinuxAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_amd64.zip"
		} else if p == platform.WindowsArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
		} else if p == platform.WindowsAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
		}

		download.Install("sqlc", url, dn)
	}

	dn := filepath.Join("lib", "database")
	yml := filepath.Join(dn, "sqlc.yaml")

	if !files.IsFile(yml) {
		messages.Fatalf("%s not found", yml)
	}

	s := spinner.New("generating queries")
	err := spinner.Start(s)
	if err != nil {
		messages.Fatal(err)
		return
	}
	sqlc := exec.Command(path.Sqlc(dn), "generate")
	sqlc.Dir = dn
	sqlc.Env = append(os.Environ())
	sqlc.Stderr = os.Stderr
	sqlc.Stdout = os.Stdout
	sqlc.Stdin = os.Stdin
	err = sqlc.Run()
	if err != nil {
		spinner.Stop(s)
		messages.Fatal(err)
	}
	spinner.Stop(s)
}
