package messages

import "fmt"

func Errorf(format string, vars ...any) {
	Error(fmt.Sprintf(format, vars...))
}
