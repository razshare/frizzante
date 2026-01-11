package messages

import "fmt"

func Sectionf(format string, vars ...any) {
	Section(fmt.Sprintf(format, vars...))
}
