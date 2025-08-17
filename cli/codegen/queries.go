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

func Queries(o QueriesOptions) error {
	if !files.IsFile(o.Sqlc) {
		var url string

		if o.Platform == platform.PlatformDarwinArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_arm64.zip"
		} else if o.Platform == platform.PlatformDarwinAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_amd64.zip"
		} else if o.Platform == platform.PlatformLinuxArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_arm64.zip"
		} else if o.Platform == platform.PlatformLinuxAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_amd64.zip"
		} else if o.Platform == platform.PlatformWindowsArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
		} else if o.Platform == platform.PlatformWindowsAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
		}

		var install Install
		install, err := Download(DownloadOptions{
			Url:  url,
			Auto: o.Auto,
		})
		if err != nil {
			return err
		}

		_, err = install(o.Sqlc)
		if err != nil {
			return err
		}
	}

	yaml := filepath.Join(o.Lib, "sqlc.yaml")

	if !files.IsFile(yaml) {
		return fmt.Errorf("%s not found", yaml)
	}

	s := spinner.New("generating queries")

	go spinner.Start(s)
	sqlc := exec.Command(o.Sqlc, "generate")
	sqlc.Dir = o.Lib
	sqlc.Env = append(os.Environ())
	sqlc.Stderr = os.Stderr
	sqlc.Stdout = os.Stdout
	sqlc.Stdin = os.Stdin
	err := sqlc.Run()
	spinner.Stop(s)

	if err != nil {
		return err
	}

	messages.Success(
		"queries generated at database.Queries.*\n",
		o.Lib+"/queries.go",
	)
	messages.Tip(
		"## usage example\n",
		"func(c *client.Client){\n",
		"    u, _ := database.Queries.FindUsers(c.Request.Context())\n",
		"    send.Json(c, u)\n",
		"}",
	)

	return nil
}
