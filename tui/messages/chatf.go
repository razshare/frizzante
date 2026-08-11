package messages

import (
	"fmt"
)

func Chatf(user string, format string, args ...any) {
	Chat(user, fmt.Sprintf(format, args...))
}
