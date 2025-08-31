package action

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/razshare/frizzante/internal/tui/messages"
)

func Dev(options DevOptions) (err error) {
	if err = Touch(TouchOptions{App: options.App}); err != nil {
		return
	}

	if err = os.MkdirAll(filepath.Join(".gen", "tmp"), os.ModePerm); err != nil {
		return
	}

	airwatch := exec.Command(options.Air)
	airwatch.Env = append(os.Environ(), "DEV=1")
	airwatch.Stderr = os.Stderr
	airwatch.Stdout = os.Stdout
	airwatch.Stdin = os.Stdin
	if err = airwatch.Start(); err != nil {
		return
	}

	messages.Success("air watcher launched")

	var group sync.WaitGroup

	group.Add(1)

	go func() { err = PackageWatch(PackageWatchOptions{App: options.App, Bun: options.Bun}) }()

	group.Wait()

	if err = airwatch.Wait(); err != nil {
		return
	}

	return
}
