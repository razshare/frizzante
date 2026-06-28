package ssr

import (
	"embed"
	"log"
)

type Options struct {
	Efs      embed.FS
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Limit    int64
}
