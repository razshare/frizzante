package messages

import "fmt"

func Infof(format string, vars ...any) {
	Info(fmt.Sprintf(format, vars...))
}
