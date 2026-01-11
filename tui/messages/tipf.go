package messages

import "fmt"

func Tipf(format string, vars ...any) {
	Tip(fmt.Sprintf(format, vars...))
}
