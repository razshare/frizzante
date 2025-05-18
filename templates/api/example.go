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

func (_ apiController) Handle(request *f.Request, response *f.Response) {
	// Noop.
}
