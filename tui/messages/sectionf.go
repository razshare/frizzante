package messages

import "fmt"

func Sectionf(format string, args ...any) {
	Section(fmt.Sprintf(format, args...))
}
