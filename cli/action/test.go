package action

import (
	"os"
	"os/exec"
)

func Test(opts TestOptions) error {
	test := exec.Command(opts.Go, "test")
	test.Env = os.Environ()
	test.Stderr = os.Stderr
	test.Stdout = os.Stdout
	test.Stdin = os.Stdin
	err := test.Run()
	if err != nil {
		return err
	}
	return nil
}
