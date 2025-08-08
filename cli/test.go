package cli

import (
	"os"
	"os/exec"
)

func OnTest() {
	OnPackage()

	test := exec.Command(Go("."), "test")
	test.Env = os.Environ()
	test.Stderr = os.Stderr
	test.Stdout = os.Stdout
	test.Stdin = os.Stdin
	runError := test.Run()
	if runError != nil {
		Fatal(runError)
	}
}
