package handlers

import (
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/sessions"
	"github.com/razshare/frizzante/views"
	"main/lib"
	"strconv"
)

func Check(con *connections.Connection) {
	state, operator := sessions.Start(con, lib.InitialState())
	defer operator.Save(state)

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

	count := int64(len(state.Todos))
	if index >= count {
		// Index is out of bounds, ignore the request.
		con.SendNavigate("/todos")
		return
	}

	state.Todos[index].Checked = true

	con.SendNavigate("/todos")
}
