package api

import (
	f "github.com/razshare/frizzante"
	//config.Server
)

var Server *f.Server
var guards []f.Guard

func init() {
	Server.OnRequest("GET /path", guards, func(req *f.Request, res *f.Response) {
		res.SendMessage("hello")
	})
}
