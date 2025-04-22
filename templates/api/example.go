package api

import f "github.com/razshare/frizzante"

func api(
	withPattern func(pattern string),
	withHandler func(handler func(req *f.Request, res *f.Response)),
) {
	withPattern("GET /")
	withHandler(func(req *f.Request, res *f.Response) {
		// Handle.
		f.SendEcho(res, "Ok.")
	})
}
