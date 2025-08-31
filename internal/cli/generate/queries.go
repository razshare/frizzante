package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/files"
	messages2 "github.com/razshare/frizzante/internal/tui/messages"
	spinner2 "github.com/razshare/frizzante/internal/tui/spinner"
)

func Queries(options QueriesOptions) (err error) {
	lib := filepath.Join("lib", "database")

	if _, err = exec.LookPath(options.Sqlc); err != nil && !files.IsFile(options.Sqlc) {
		if err = Sqlc(SqlcOptions{Sqlc: options.Sqlc, Platform: options.Platform, Auto: options.Auto}); err != nil {
			return
		}
	}

	yaml := filepath.Join(lib, "sqlc.yaml")

	if !files.IsFile(yaml) {
		return fmt.Errorf("%s not found", yaml)
	}

	spin := spinner2.New("generating queries")

	var sqlc string
	if files.IsFile(options.Sqlc) {
		if sqlc, err = filepath.Rel(lib, options.Sqlc); err != nil {
			return err
		}
	} else if sqlc, err = exec.LookPath(options.Sqlc); err != nil {
		sqlc = options.Sqlc
	}

	go spinner2.Start(spin)
	generate := exec.Command(sqlc, "generate")
	generate.Dir = lib
	generate.Env = append(os.Environ())
	generate.Stderr = os.Stderr
	generate.Stdout = os.Stdout
	generate.Stdin = os.Stdin
	err = generate.Run()
	spinner2.Stop(spin)

	if err != nil {
		return
	}

	messages2.Success(
		"queries generated at database.Queries.*\n",
		lib+"/queries.go",
	)
	messages2.Tip(
		"## usage example\n",
		"func(c *client.Client){\n",
		"    u, _ := database.Queries.FindUsers(c.Request.Context())\n",
		"    send.Json(c, u)\n",
		"}",
	)

	return
}
