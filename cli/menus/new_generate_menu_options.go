package menus

import (
	"embed"

	"github.com/razshare/frizzante/cli/apps"
)

type NewGenerateMenuOptions struct {
	Efs        embed.FS
	Modifiers  apps.Modifiers
	Persistent bool
}
