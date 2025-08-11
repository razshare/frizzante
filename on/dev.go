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

	mkdirError := os.MkdirAll(filepath.Join(".gen", "tmp"), os.ModePerm)
	if mkdirError != nil {
		messages.Fatal(mkdirError)
	}

	air := exec.Command(path.Air("."))
	air.Env = append(os.Environ(), "DEV=1")
	air.Stderr = os.Stderr
	air.Stdout = os.Stdout
	air.Stdin = os.Stdin
	airError := air.Start()
	if airError != nil {
		messages.Fatalf("air watcher failed to launch\n%s", airError)
	}
	messages.Success("air watcher launched")

	var group sync.WaitGroup

	group.Add(1)

	go func() { PackageWatch() }()

	group.Wait()
	tidyWaitError := air.Wait()
	if tidyWaitError != nil {
		messages.Fatal(tidyWaitError)
	}
}
