package messages

import "fmt"

func Successf(format string, vars ...any) {
	Success(fmt.Sprintf(format, vars...))
}
