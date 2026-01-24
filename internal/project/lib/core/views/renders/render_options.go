package renders

import (
	"embed"
	"log"

	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

type RenderOptions struct {
	View     views.View
	Efs      embed.FS
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
