package action

import (
	"fmt"
	"github.com/razshare/frizzante/files"
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

	var air string
	if files.IsFile(opts.Air) {
		if air, err = filepath.Rel(opts.App, opts.Air); err != nil {
			return err
		}
	} else if air, err = exec.LookPath(opts.Air); err != nil {
		air = opts.Air
	}

	airwatch := exec.Command(air)
	airwatch.Env = append(os.Environ(), "DEV=1")
	airwatch.Stderr = os.Stderr
	airwatch.Stdout = os.Stdout
	airwatch.Stdin = os.Stdin
	err = airwatch.Start()
	if err != nil {
		return fmt.Errorf("air watcher failed to launch\n%s", err)
	}

	messages.Success("air watcher launched")

	var group sync.WaitGroup

	group.Add(1)

	go func() { err = PkgWatch(PkgWatchOptions{App: opts.App, Bun: opts.Bun}) }()

	group.Wait()

	err = airwatch.Wait()
	if err != nil {
		return
	}

	return nil
}
