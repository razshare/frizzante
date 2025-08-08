package cli

import (
	"os"
	"os/exec"
)

func OnFormat() {
	OnTouch()

	gofmt := exec.Command(Go("."), "fmt")
	gofmt.Env = append(os.Environ())
	gofmt.Stderr = os.Stderr
	gofmt.Stdout = os.Stdout
	gofmt.Stdin = os.Stdin
	gofmtError := gofmt.Run()
	if gofmtError != nil {
		Fatal(gofmtError)
	}

	prettier := exec.Command(Bun("app"), "x", "prettier", "--write", ".")
	prettier.Dir = "app"
	prettier.Env = append(os.Environ())
	prettier.Stderr = os.Stderr
	prettier.Stdout = os.Stdout
	prettier.Stdin = os.Stdin
	prettierError := prettier.Run()
	if prettierError != nil {
		Fatal(prettierError)
	}

	Success("project formatted")
}
