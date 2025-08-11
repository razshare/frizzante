package on

import (
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
)

func Test() {
	Package()

	test := exec.Command(path.Go("."), "test")
	test.Env = os.Environ()
	test.Stderr = os.Stderr
	test.Stdout = os.Stdout
	test.Stdin = os.Stdin
	err := test.Run()
	if err != nil {
		messages.Fatal(err)
	}
}
