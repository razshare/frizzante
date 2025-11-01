package render_function

import (
	"log"

	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

type LogLevel uint8

const LogLevelBase LogLevel = 0
const LogLevelWarning LogLevel = 1
const LogLevelDanger LogLevel = 2

type Config struct {
	Data     []byte
	App      string
	Server   string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

type RenderFunction = func(view views.View) (head string, body string, err error)
