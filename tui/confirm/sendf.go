package confirm

import (
	"fmt"
)

func Sendf(defaultValue bool, format string, vars ...any) (yes bool, err error) {
	return Send(defaultValue, fmt.Sprintf(format, vars...))
}
