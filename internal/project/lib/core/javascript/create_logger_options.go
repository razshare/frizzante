package javascript

import (
	"log"

	"github.com/dop251/goja"
)

type CreateLoggerOptions struct {
	Level    LogLevel
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Runtime  *goja.Runtime
}
