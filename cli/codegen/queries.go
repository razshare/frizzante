package codegen

import (
	"fmt"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/platform"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"path/filepath"
)

func Queries(clr bool, sqlcbin string) error {
	if !files.IsFile(sqlcbin) {
		dst := filepath.Join(".gen", "sqlc")

		var plat platform.Type
		plat, err := platform.Find(clr)

		if err != nil {
			return err
		}

		var url string

		if plat == platform.DarwinArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_arm64.zip"
		} else if plat == platform.DarwinAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_amd64.zip"
		} else if plat == platform.LinuxArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_arm64.zip"
		} else if plat == platform.LinuxAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_amd64.zip"
		} else if plat == platform.WindowsArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
		} else if plat == platform.WindowsAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
		}

		var install Install
		install, err = Download(url)
		if err != nil {
			return err
		}

		_, err = install(dst)
		if err != nil {
			return err
		}
	}

	database := filepath.Join("lib", "database")
	sqlcyaml := filepath.Join(database, "sqlc.yaml")

	if !files.IsFile(sqlcyaml) {
		return fmt.Errorf("%s not found", sqlcyaml)
	}

	s := spinner.New("generating queries")

	sqlcbin, err := path.Sqlc(sqlcbin, database)
	if err != nil {
		return err
	}

	go spinner.Start(s)
	sqlc := exec.Command(sqlcbin, "generate")
	sqlc.Dir = database
	sqlc.Env = append(os.Environ())
	sqlc.Stderr = os.Stderr
	sqlc.Stdout = os.Stdout
	sqlc.Stdin = os.Stdin
	err = sqlc.Run()
	spinner.Stop(s)

	if err != nil {
		return err
	}

	messages.Success(
		"queries generated at database.Queries.*\n",
		database+"/queries.go",
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
