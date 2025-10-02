package render

import (
	"embed"
	"log"
)

type Config struct {
	App      string
	Efs      embed.FS
	Limit    int
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
