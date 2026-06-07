package queries

import (
	"os"
	"os/exec"
	"path/filepath"
)

func Generate() (err error) {
	command := exec.Command(filepath.Join(".gen", "sqlc", "sqlc"), "generate")
	command.Env = os.Environ()
	if err = command.Run(); err != nil {
		return
	}
	return
}
