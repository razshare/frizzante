package generate

import (
	"path/filepath"

	"github.com/razshare/frizzante/cli/npm"
	"github.com/razshare/frizzante/tui/messages"
)

func Icons(options IconsOptions) (err error) {
	if err = Copy(CopyOptions{
		From: "internal/additions/app/lib/components/icons",
		To:   filepath.Join(options.App, "lib", "components", "icons"),
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if err = npm.Install(options.Bun, options.App, "@mdi/js"); err != nil {
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
