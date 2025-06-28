package handlers

import (
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/sessions"
	"github.com/razshare/frizzante/views"
	"main/lib"
	"strconv"
)

func Remove(con *connections.Connection) {
	state, operator := sessions.Start(con, lib.InitialState())
	defer operator.Save(state)

	count := int64(len(state.Todos))

	if 0 == count {
		// No index found, ignore the request.
		con.SendNavigate("/todos")
		return
	}

	indexString := con.ReceiveQuery("index")
	if "" == indexString {
		// No index found, ignore the request.
		con.SendNavigate("/todos")
		return
	}

	index, indexError := strconv.ParseInt(indexString, 10, 64)
	if nil != indexError {
		con.SendView(views.View{Name: "Todos", Data: map[string]any{
			"error": indexError.Error(),
		}})
		return
	}

	if index >= count {
		// Index is out of bounds, ignore the request.
		con.SendNavigate("/todos")
		return
	}

	state.Todos = append(state.Todos[:index], state.Todos[index+1:]...)

	con.SendNavigate("/todos")
}
