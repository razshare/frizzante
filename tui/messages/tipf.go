package messages

import "fmt"

func Tipf(format string, args ...any) {
	Tip(fmt.Sprintf(format, args...))
}
