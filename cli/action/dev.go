package action

import (
	"fmt"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func Dev(opts DevOptions) (err error) {
	err = Touch(TouchOptions{App: opts.App})
	if err != nil {
		return
	}

	err = os.MkdirAll(filepath.Join(".gen", "tmp"), os.ModePerm)
	if err != nil {
		return
	}

	air := exec.Command(opts.Air)
	air.Env = append(os.Environ(), "DEV=1")
	air.Stderr = os.Stderr
	air.Stdout = os.Stdout
	air.Stdin = os.Stdin
	err = air.Start()
	if err != nil {
		return fmt.Errorf("air watcher failed to launch\n%s", err)
	}

	messages.Success("air watcher launched")

	var group sync.WaitGroup

	group.Add(1)

	go func() { err = PkgWatch(PkgWatchOptions{App: opts.App, Bun: opts.Bun}) }()

	group.Wait()

	err = air.Wait()
	if err != nil {
		return
	}

	return nil
}
