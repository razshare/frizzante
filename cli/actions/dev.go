package actions

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/razshare/frizzante/cli/generations"
	tags_ "github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/tui/messages"
)

func Dev(options DevOptions) (err error) {
	if err = os.MkdirAll(filepath.Join("app", "dist"), os.ModePerm); err != nil {
		return
	}
	if err = os.MkdirAll(filepath.Join(".gen", "tmp"), os.ModePerm); err != nil {
		return
	}
	var tags []string
	if tags, err = tags_.Parse(options.Tags); err != nil {
		return
	}
	if !slices.Contains(tags, "dev") {
		tags = append(tags, "dev")
	}
	if !slices.Contains(tags, "trace") {
		tags = append(tags, "trace")
	}
	if err = generations.AirConfig(generations.AirConfigOptions{
		Efs:  options.Efs,
		Tags: strings.Join(tags, ","),
	}); err != nil {
		return
	}
	var group sync.WaitGroup
	group.Go(func() {
		if perr := PackageWatch(PackageWatchOptions{
			Bun: options.Bun,
		}); perr != nil {
			messages.Error(perr)
		}
	})
	group.Go(func() {
		if !messages.Command(messages.CommandOptions{
			Environment: os.Environ(),
			Program:     options.Air,
		}) {
			messages.Error("air failed")
		}
	})
	group.Wait()
	return
}
