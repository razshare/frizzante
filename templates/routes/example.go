package routes

import "github.com/razshare/frizzante"

func functionName(req *frizzante.Request, res *frizzante.Response) {
	res.SendView(frizzante.View{Name: "viewName"})
}
