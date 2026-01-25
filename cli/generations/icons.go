package generations

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/cli/npm"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
)

func Icons(options IconsOptions) (err error) {
	directoryName := filepath.Join("app", "lib", "components", "icons")
	if files.IsDirectory(directoryName) {
		if options.Strict {
			err = fmt.Errorf("%s already exists", directoryName)
			return
		}
		var yesRemove bool
		if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", directoryName); err != nil {
			return err
		}
		if yesRemove {
			if err = os.RemoveAll(directoryName); err != nil {
				return
			}
		}
	}
	if err = Copy(CopyOptions{
		From: "internal/project/app/lib/components/icons",
		To:   filepath.Join("app", "lib", "components", "icons"),
		Efs:  options.Efs,
	}); err != nil {
		return
	}
	if err = npm.Install(npm.InstallOptions{
		Bun:          options.Bun,
		PackageNames: []string{"@mdi/js"},
	}); err != nil {
		return
	}
	messages.Tip(
		"## usage example\n",
		"<script lang=\"ts\">\n",
		"    import Icon from '$lib/components/icons/Icon.svelte'\n",
		"    import { mdiClose } from '@mdi/js'\n",
		"</script>\n",
		"\n",
		"<Icon path={mdiClose}/>",
	)
	return
}
