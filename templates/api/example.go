package api

import f "github.com/razshare/frizzante"

var apiName = f.
	NewApiController().
	WithPattern("GET /path").
	WithHandler(handler)

func handler(req *f.Request, res *f.Response) {
	res.SendMessage("hello")
}
