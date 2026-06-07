package messages

import "fmt"

func Errorf(format string, args ...any) {
	Error(fmt.Sprintf(format, args...))
}
