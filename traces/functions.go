package traces

import (
	"fmt"
	"log"
	"runtime"
)

func Trace(logger *log.Logger, value any) {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		logger.Printf("%s:%d %s\n", file, line, value)
	} else {
		logger.Println(value)
	}
}

func Tracef(logger *log.Logger, format string, values ...any) {
	Trace(logger, fmt.Sprintf(format, values...))
}
