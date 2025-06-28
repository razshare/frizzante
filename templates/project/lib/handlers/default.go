package handlers

import "github.com/razshare/frizzante/connections"

func Default(con *connections.Connection) {
	con.SendFileOrElse(func() { Welcome(con) })
}
