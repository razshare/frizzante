package render

import (
	"embed"
	"log"

	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

type FunctionOptions struct {
	View     views.View
	Efs      embed.FS
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
