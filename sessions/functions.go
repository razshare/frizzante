package sessions

import (
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/operators"
	"path/filepath"
)

var archiveGlobal archives.Archive

func init() {
	diskArchive := archives.NewDiskArchive()
	diskArchive.Name = filepath.Join(".gen", "sessions")
	archiveGlobal = diskArchive
}

// WithArchive sets the sessions archive.
func WithArchive(archive archives.Archive) {
	archiveGlobal = archive
}

// Start starts a session.
func Start[T any](con *connections.Connection, state T) (*T, *operators.Operator) {
	operator := &operators.Operator{
		Archive:    archiveGlobal,
		Connection: con,
	}

	if !operator.Exists() {
		operator.Save(&state)
	} else {
		operator.Load(&state)
	}

	return &state, operator
}

// StartEmpty starts a session with empty state.
func StartEmpty[T any](con *connections.Connection) (*T, *operators.Operator) {
	var state T
	return Start(con, state)
}
