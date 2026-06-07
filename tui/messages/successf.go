package messages

import "fmt"

func Successf(format string, args ...any) {
	Success(fmt.Sprintf(format, args...))
}
