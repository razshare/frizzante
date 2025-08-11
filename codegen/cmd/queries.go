package cmd

import (
	"embed"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"path/filepath"
)

func Queries(efs embed.FS) {
	if !files.IsFile(filepath.Join(".gen", "sqlc", "sqlc")) {
		messages.Fatal("sqlc not found")
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
