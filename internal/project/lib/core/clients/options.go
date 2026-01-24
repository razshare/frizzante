package clients

import (
	"embed"
	"log"

	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

type Options struct {
	Efs      embed.FS
	Render   renders.Render
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
