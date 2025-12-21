package apps

import "embed"

type App struct {
	Modifiers Modifiers
	Efs       embed.FS
}
