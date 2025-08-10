package cli

import (
	"os"
	"os/exec"
)

func OnSqlcGenerate() {
	sqlcGenerate := exec.Command(Sqlc("."), "generate")
	sqlcGenerate.Env = append(os.Environ())
	sqlcGenerate.Stderr = os.Stderr
	sqlcGenerate.Stdout = os.Stdout
	sqlcGenerate.Stdin = os.Stdin
	svelteCheckError := sqlcGenerate.Run()
	if svelteCheckError != nil {
		Fatal(svelteCheckError)
	}
	Success("files generated")
}
