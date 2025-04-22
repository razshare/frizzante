package api

import f "github.com/razshare/frizzante"

func api(
	withPattern func(pattern string),
	withHandler func(handler func(
		request *f.Request,
		response *f.Response,
	)),
) {
	withPattern("GET /")
	withHandler(func(
		request *f.Request,
		response *f.Response,
	) {
		// Handle.
		f.SendEcho(response, "Ok.")
	})
}
