package on

import (
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func Dev() {
	Touch()

	err := os.MkdirAll(filepath.Join(".gen", "tmp"), os.ModePerm)
	if err != nil {
		messages.Fatal(err)
	}

	air := exec.Command(path.Air("."))
	air.Env = append(os.Environ(), "DEV=1")
	air.Stderr = os.Stderr
	air.Stdout = os.Stdout
	air.Stdin = os.Stdin
	err = air.Start()
	if err != nil {
		messages.Fatalf("air watcher failed to launch\n%s", err)
	}

	messages.Success("air watcher launched")

	var group sync.WaitGroup

	group.Add(1)

	go func() { PackageWatch() }()

	group.Wait()

	err = air.Wait()
	if err != nil {
		messages.Fatal(err)
	}
}
