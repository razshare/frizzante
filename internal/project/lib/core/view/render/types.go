package render

import (
	"embed"
	"log"

	view_ "github.com/razshare/frizzante/internal/project/lib/core/view"
)

type Config struct {
	App      string
	Efs      embed.FS
	Limit    int
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

type Render = func(view view_.View) (html string, err error)
