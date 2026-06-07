package logs

import (
	"fmt"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

func Infof(client *clients.Client, format string, args ...any) {
	Info(client, fmt.Sprintf(format, args...))
}
