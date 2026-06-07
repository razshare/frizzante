package messages

import "fmt"

func Fatalf(format string, args ...any) {
	Fatal(fmt.Sprintf(format, args...))
}
