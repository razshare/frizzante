package cli

import (
	"os"
	"os/exec"
)

func OnUpdate() {
	OnTouch()

	get := exec.Command(Go("."), "get", "-u", "./...")
	get.Env = append(os.Environ())
	get.Stderr = os.Stderr
	get.Stdout = os.Stdout
	get.Stdin = os.Stdin
	getError := get.Run()
	if getError != nil {
		Fatal(getError)
	}

	prettier := exec.Command(Bun("app"), "update")
	prettier.Dir = "app"
	prettier.Env = append(os.Environ())
	prettier.Stderr = os.Stderr
	prettier.Stdout = os.Stdout
	prettier.Stdin = os.Stdin
	prettierError := prettier.Run()
	if prettierError != nil {
		Fatal(prettierError)
	}

	Success("project dependencies updated")
}
