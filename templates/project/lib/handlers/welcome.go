package handlers

import (
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/views"
)

func Welcome(con *connections.Connection) {
	con.SendView(views.View{Name: "Welcome"})
}
