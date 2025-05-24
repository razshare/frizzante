package api

import f "github.com/razshare/frizzante"

type Controller struct{}

func (_ Controller) Configure() f.ApiConfiguration {
	return f.ApiConfiguration{
		Pattern: "GET /path",
	}
}

func (_ Controller) Handle(req *f.Request, res *f.Response) {
	res.SendMessage("Hello.")
}
