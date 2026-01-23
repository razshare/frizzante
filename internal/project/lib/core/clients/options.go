package clients

import (
	"embed"
	"log"

	"github.com/razshare/frizzante/internal/project/lib/core/views/render"
)

type Options struct {
	Efs      embed.FS
	Render   render.Function
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
