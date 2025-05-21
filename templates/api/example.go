package api

import f "github.com/razshare/frizzante"

var apiName = f.
	NewApiController().
	WithPath("/path").
	WithHandler("GET", apiGet)

func apiGet(req *f.Request, res *f.Response) {
	res.SendMessage("hello")
}
