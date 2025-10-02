package compile_

import (
	"log"
)

type LogLevel uint8

const LogLevelBase LogLevel = 0
const LogLevelWarning LogLevel = 1
const LogLevelDanger LogLevel = 2

type Config struct {
	Data     []byte
	Format   string
	App      string
	Server   string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
