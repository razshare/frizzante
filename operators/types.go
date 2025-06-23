package operators

import (
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/connections"
)

type Operator struct {
	Archive    archives.Archive
	Connection *connections.Connection
}
