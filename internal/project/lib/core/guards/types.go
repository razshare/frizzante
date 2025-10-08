package guards

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

type Guard struct {
	Name    string
	Handler func(client *clients.Client, allow func())
}
