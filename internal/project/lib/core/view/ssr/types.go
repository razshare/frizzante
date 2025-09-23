package ssr

import (
	"embed"
	"log"
)

type LogLevel uint8

const LogLevelBase LogLevel = 0
const LogLevelWarning LogLevel = 1
const LogLevelDanger LogLevel = 2

type Config struct {
	App      string
	Efs      embed.FS
	Limit    int
	Disk     bool
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
