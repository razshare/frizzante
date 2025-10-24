package actions

import (
	"os"

	"github.com/razshare/frizzante/cli/detect"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Reset(_ ResetOptions) (err error) {
	var cache string
	if cache, err = detect.FrizzanteCache(); err != nil {
		return
	}

	if files.IsDirectory(cache) {
		if err = os.RemoveAll(cache); err != nil {
			return
		}

		messages.Successf("%s deleted", cache)

		return
	}

	messages.Infof("%s not found", cache)

	return
}
