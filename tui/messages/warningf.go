package messages

import "fmt"

func Warningf(format string, vars ...any) {
	Warning(fmt.Sprintf(format, vars...))
}
