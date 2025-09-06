package action

import (
	"os"
	"os/exec"
)

func Test(options TestOptions) (err error) {
	test := exec.Command(options.Go, "test", "./...")
	test.Env = os.Environ()
	test.Stderr = os.Stderr
	test.Stdout = os.Stdout
	test.Stdin = os.Stdin
	err = test.Run()

	return
}
