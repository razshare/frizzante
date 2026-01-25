package inputs

import "fmt"

func Sendf(format string, vars ...any) (value string, err error) {
	return Send(fmt.Sprintf(format, vars...))
}
