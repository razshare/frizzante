package actions

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/tui/messages"
)

func Dev(options DevOptions) (err error) {
	if err = Touch(TouchOptions{}); err != nil {
		return
	}

	if err = os.MkdirAll(filepath.Join(".gen", "tmp"), os.ModePerm); err != nil {
		return
	}

	if err = generate.AirConfig(generate.AirConfigOptions{
		Efs:  options.Efs,
		Tags: options.Tags,
	}); err != nil {
		return
	}

	var group sync.WaitGroup
	group.Go(func() {
		_ = PackageWatch(PackageWatchOptions{Bun: options.Bun})
	})
	group.Go(func() {
		messages.Command(".", os.Environ(), options.Air)
	})
	group.Wait()

	return
}
