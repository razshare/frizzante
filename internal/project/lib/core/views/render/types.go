package render

import (
	"embed"
	"log"

	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

type Render = func(options Options) (html string, err error)

type Options struct {
	View     views.View
	Efs      embed.FS
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
