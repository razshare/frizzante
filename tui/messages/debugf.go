package messages

import "fmt"

func Debugf(format string, vars ...any) {
	Debug(fmt.Sprintf(format, vars...))
}
