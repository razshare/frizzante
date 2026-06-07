package messages

import "fmt"

func Debugf(format string, args ...any) {
	Debug(fmt.Sprintf(format, args...))
}
