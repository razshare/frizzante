package api

import f "github.com/razshare/frizzante"

func handler(req *f.Request, res *f.Response) {
	// Serve api.
}

func api(
	withPattern func(string),
	withHandler func(handler func(req *f.Request, res *f.Response)),
) {
	withPattern("GET /")
	withHandler(handler)
}
