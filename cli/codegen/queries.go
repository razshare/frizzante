package codegen

import (
	"fmt"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"path/filepath"
)

func Queries(opts QueriesOptions) (err error) {
	if _, perr := exec.LookPath(opts.Sqlc); perr != nil && !files.IsFile(opts.Sqlc) {
		if err = Sqlc(SqlcOptions{Sqlc: opts.Sqlc, Platform: opts.Platform, Auto: opts.Auto}); err != nil {
			return
		}
	}

	yaml := filepath.Join(opts.Lib, "sqlc.yaml")

	if !files.IsFile(yaml) {
		return fmt.Errorf("%s not found", yaml)
	}

	spin := spinner.New("generating queries")

	var sqlc string
	if files.IsFile(opts.Sqlc) {
		if sqlc, err = filepath.Rel(opts.Lib, opts.Sqlc); err != nil {
			return err
		}
	} else if sqlc, err = exec.LookPath(opts.Sqlc); err != nil {
		sqlc = opts.Sqlc
	}

	go spinner.Start(spin)
	gen := exec.Command(sqlc, "generate")
	gen.Dir = opts.Lib
	gen.Env = append(os.Environ())
	gen.Stderr = os.Stderr
	gen.Stdout = os.Stdout
	gen.Stdin = os.Stdin
	err = gen.Run()
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
