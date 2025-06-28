package handlers

import (
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/sessions"
	"github.com/razshare/frizzante/views"
	"main/lib"
)

func Todos(con *connections.Connection) {
	state, operator := sessions.Start(con, lib.InitialState())
	defer operator.Save(state)

	con.SendView(views.View{Name: "Todos", Data: map[string]any{
		"todos": state.Todos,
	}})
}
