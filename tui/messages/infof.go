package messages

import "fmt"

func Infof(format string, args ...any) {
	Info(fmt.Sprintf(format, args...))
}
