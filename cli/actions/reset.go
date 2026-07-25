package actions

import (
	"os"

	"github.com/razshare/frizzante/v2/cli/paths"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/tui/messages"
)

func Reset(_ ResetOptions) (err error) {
	var cache string
	if cache, err = paths.Cache(); err != nil {
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
