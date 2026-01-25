package spinners

import "fmt"

func Newf(format string, vars ...any) *Spinner {
	return New(fmt.Sprintf(format, vars...))
}
