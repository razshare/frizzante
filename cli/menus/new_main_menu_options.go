package menus

import (
	"embed"

	"github.com/razshare/frizzante/cli/apps"
)

type NewMainMenuOptions struct {
	Efs        embed.FS
	Modifiers  apps.Modifiers
	Persistent bool
}
