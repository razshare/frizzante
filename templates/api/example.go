package api

import (
	f "github.com/razshare/frizzante"
)

type apiController struct {
	f.ApiController
}

func (_ apiController) Configure() f.ApiConfiguration {
	return f.ApiConfiguration{
		Pattern: "GET /path",
	}
}

func (_ apiController) Handle(req *f.Request, res *f.Response) {
	res.SendMessage("hello")
}
