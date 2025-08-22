package codegen

import (
	"fmt"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/platform"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"path/filepath"
)

func Queries(opts QueriesOptions) (err error) {
	if !files.IsFile(opts.Sqlc) {
		var url string

		if opts.Platform == platform.DarwinArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_arm64.zip"
		} else if opts.Platform == platform.DarwinAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_amd64.zip"
		} else if opts.Platform == platform.LinuxArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_arm64.zip"
		} else if opts.Platform == platform.LinuxAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_amd64.zip"
		} else if opts.Platform == platform.WindowsArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
		} else if opts.Platform == platform.WindowsAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
		}

		var install Install
		if install, _, err = Download(DownloadOptions{Url: url, Auto: opts.Auto}); err != nil {
			return
		}

		if _, err = install(opts.Sqlc); err != nil {
			return
		}
	}

	yaml := filepath.Join(opts.Lib, "sqlc.yaml")

	if !files.IsFile(yaml) {
		return fmt.Errorf("%s not found", yaml)
	}

	spin := spinner.New("generating queries")

	go spinner.Start(spin)
	sqlc := exec.Command(opts.Sqlc, "generate")
	sqlc.Dir = opts.Lib
	sqlc.Env = append(os.Environ())
	sqlc.Stderr = os.Stderr
	sqlc.Stdout = os.Stdout
	sqlc.Stdin = os.Stdin
	err = sqlc.Run()
	spinner.Stop(spin)

	if err != nil {
		return
	}

	messages.Success(
		"queries generated at database.Queries.*\n",
		opts.Lib+"/queries.go",
	)
	messages.Tip(
		"## usage example\n",
		"func(c *client.Client){\n",
		"    u, _ := database.Queries.FindUsers(c.Request.Context())\n",
		"    send.Json(c, u)\n",
		"}",
	)

	return
}
