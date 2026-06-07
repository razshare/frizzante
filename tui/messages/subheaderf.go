package messages

import "fmt"

func Subheaderf(format string, args ...any) {
	Subheader(fmt.Sprintf(format, args...))
}
