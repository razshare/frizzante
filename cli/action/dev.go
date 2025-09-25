package action

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/razshare/frizzante/tui/messages"
)

func Dev(options DevOptions) (err error) {
	if err = Touch(TouchOptions{App: options.App}); err != nil {
		return
	}

	if err = os.MkdirAll(filepath.Join(".gen", "tmp"), os.ModePerm); err != nil {
		return
	}

	var group sync.WaitGroup
	group.Go(func() {
		messages.Command(".", append(os.Environ(), "DEV=1"), options.Air)
	})
	group.Go(func() {
		_ = PackageWatch(PackageWatchOptions{App: options.App, Bun: options.Bun})
	})
	group.Wait()

	return
}
