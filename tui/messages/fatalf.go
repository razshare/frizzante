package messages

import "fmt"

func Fatalf(format string, vars ...any) {
	Fatal(fmt.Sprintf(format, vars...))
}
