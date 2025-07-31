package guards

import (
	"github.com/razshare/frizzante/connections"
)

type Guard struct {
	Name    string
	Handler func(connection *connections.Connection, allow func())
	Tags    []string
}
