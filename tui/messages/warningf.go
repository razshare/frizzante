package messages

import "fmt"

func Warningf(format string, args ...any) {
	Warning(fmt.Sprintf(format, args...))
}
