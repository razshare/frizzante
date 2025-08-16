package on

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/path"
	"os"
	"os/exec"
)

func Test(c *cli.Cli) error {
	err := Package(c)
	if err != nil {
		return err
	}

	gobin, err := path.Go(c, ".")
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
