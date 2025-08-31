package action

import (
	"os"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/internal/cli/user"
	"github.com/razshare/frizzante/internal/tui/messages"
)

func Reset(_ ResetOptions) (err error) {
	var cache string
	if cache, err = user.FrizzanteCache(); err != nil {
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
