package on

import (
	"os"
	"os/exec"
)

func Test(app string, gobin string, bunbin string) error {
	err := Package(app, bunbin)
	if err != nil {
		return err
	}

	test := exec.Command(gobin, "test")
	test.Env = os.Environ()
	test.Stderr = os.Stderr
	test.Stdout = os.Stdout
	test.Stdin = os.Stdin
	err = test.Run()
	if err != nil {
		return err
	}
	return nil
}
