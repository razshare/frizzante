package action

import (
	"os"
	"os/exec"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Test(options TestOptions) (err error) {
	spin := spinner.New("testing code")
	go spinner.Start(spin)
	defer spinner.Stop(spin)

	test := exec.Command(options.Go, "test", "./...")
	test.Env = os.Environ()
	test.Stderr = os.Stderr
	test.Stdout = os.Stdout
	test.Stdin = os.Stdin
	err = test.Run()

	if test.Err != nil {
		messages.Error(test.Err.Error())
	}
	return
}
