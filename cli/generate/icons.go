package generate

import (
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/cli/npm"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Icons(options IconsOptions) (err error) {
	if err = Copy(CopyOptions{
		From: "internal/project/app/lib/components/icons",
		To:   filepath.Join("app", "lib", "components", "icons"),
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel("app", options.Bun); err != nil {
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	if bun == "" {
		err = errors.New("bun not found")
		return
	}

	if err = npm.Install(npm.InstallOptions{
		Bun:      bun,
		Packages: []string{"@mdi/js"},
	}); err != nil {
		return
	}

	messages.Tip(
		"## usage example\n",
		"<script lang=\"ts\">\n",
		"    import Icon from '@lib/components/icons/Icon.svelte'\n",
		"    import { mdiClose } from '@mdi/js'\n",
		"</script>\n",
		"\n",
		"<Icon path={mdiClose}/>",
	)

	return
}
