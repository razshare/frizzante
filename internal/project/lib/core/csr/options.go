package csr

import (
	"embed"
	"log"
)

type Options struct {
	Efs      embed.FS
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
