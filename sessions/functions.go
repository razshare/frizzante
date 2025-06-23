package sessions

import (
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/operators"
	"path/filepath"
)

// Start starts a session.
func Start[T any](con *connections.Connection) (*T, *operators.Operator) {
	var state T
	archive := archives.NewDiskArchive()
	archive.Name = filepath.Join(".gen", "sessions")
	operator := &operators.Operator{
		Archive:    archive,
		Connection: con,
	}

	if !operator.Exists() {
		operator.Save(&state)
	} else {
		operator.Load(&state)
	}

	return &state, operator
}
