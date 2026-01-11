package messages

import "fmt"

func Subheaderf(format string, vars ...any) {
	Subheader(fmt.Sprintf(format, vars...))
}
