package api

import (
	f "github.com/razshare/frizzante"
	//config.Server
)

var Server *f.Server

func init() {
	Server.OnRequest("GET /path", func(req *f.Request, res *f.Response) {
		res.SendMessage("hello")
	})
}
