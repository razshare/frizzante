package generate

import (
	"fmt"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"path/filepath"
)

func Queries(options QueriesOptions) (err error) {
	if _, err = exec.LookPath(options.Sqlc); err != nil && !files.IsFile(options.Sqlc) {
		if err = Sqlc(SqlcOptions{Sqlc: options.Sqlc, Platform: options.Platform, Auto: options.Auto}); err != nil {
			return
		}
	}

	yaml := filepath.Join(options.Lib, "sqlc.yaml")

	if !files.IsFile(yaml) {
		return fmt.Errorf("%s not found", yaml)
	}

	spin := spinner.New("generating queries")

	var sqlc string
	if files.IsFile(options.Sqlc) {
		if sqlc, err = filepath.Rel(options.Lib, options.Sqlc); err != nil {
			return err
		}
	} else if sqlc, err = exec.LookPath(options.Sqlc); err != nil {
		sqlc = options.Sqlc
	}

	go spinner.Start(spin)
	generate := exec.Command(sqlc, "generate")
	generate.Dir = options.Lib
	generate.Env = append(os.Environ())
	generate.Stderr = os.Stderr
	generate.Stdout = os.Stdout
	generate.Stdin = os.Stdin
	err = generate.Run()
	spinner.Stop(spin)

	if err != nil {
		return
	}

	messages.Success(
		"queries generated at database.Queries.*\n",
		options.Lib+"/queries.go",
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
