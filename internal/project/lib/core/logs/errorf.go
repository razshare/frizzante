package logs

import (
	"fmt"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

func Errorf(client *clients.Client, format string, args ...any) {
	Error(client, fmt.Sprintf(format, args...))
}
