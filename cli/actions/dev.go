package actions

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/razshare/frizzante/cli/generate"
	tags_ "github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
)

func Dev(options DevOptions) (err error) {
	if err = Touch(TouchOptions{}); err != nil {
		return
	}

	if err = os.MkdirAll(filepath.Join(".gen", "tmp"), os.ModePerm); err != nil {
		return
	}

	var tags []string
	if tags, err = tags_.Parse(options.Tags); err != nil {
		return
	}

	if len(tags) == 0 && options.Interactive {
		if tags, err = tags_.Select([]search.Choice{
			{Id: "types", Description: "enables type generations"},
			{Id: "no_js_runtime", Description: "disables the server-side JavaScript runtime"},
			{Id: "experimental_qjs_runtime", Description: "replaces goja with qjs"},
			{Id: "other", Description: "adds custom tags"},
		}); err != nil {
			return
		}
	}

	if err = generate.AirConfig(generate.AirConfigOptions{
		Efs:  options.Efs,
		Tags: strings.Join(append(tags, "dev", "trace"), ","),
	}); err != nil {
		return
	}

	var group sync.WaitGroup
	group.Go(func() {
		_ = PackageWatch(PackageWatchOptions{Bun: options.Bun})
	})
	group.Go(func() {
		messages.Command(messages.CommandOptions{Environment: os.Environ(), Program: options.Air})
	})
	group.Wait()

	return
}
