package api

import f "github.com/razshare/frizzante"

func api(
	withPattern func(string),
	withHandler func(func(*f.Request, *f.Response)),
) {
	withPattern("GET /")
	withHandler(func(request *f.Request, response *f.Response) {
		// Handle.
		f.SendEcho(response, "Ok.")
	})
}
