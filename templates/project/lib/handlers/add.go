package handlers

import (
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/sessions"
	"github.com/razshare/frizzante/views"
	"main/lib"
)

func Add(con *connections.Connection) {
	state, operator := sessions.Start(con, lib.InitialState())
	defer operator.Save(state)

	description := con.ReceiveQuery("description")

	if "" == description {
		con.SendView(views.View{Name: "Todos", Data: map[string]any{
			"todos": state.Todos,
			"error": "todo description cannot be empty",
		}})
		return
	}

	state.Todos = append(state.Todos, lib.Todo{Checked: false, Description: description})

	con.SendNavigate("/todos")
}
