package sessions

import (
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/operators"
)

type SessionStarter[T any] = func(con *connections.Connection, state T) (*T, *operators.Operator)
