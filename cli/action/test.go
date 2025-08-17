package action

import (
	"os"
	"os/exec"
)

func Test(o TestOptions) error {
	err := Pkg(PkgOptions{App: o.App, Bun: o.Bun})
	if err != nil {
		return err
	}

	test := exec.Command(o.Go, "test")
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
