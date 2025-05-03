package api

import f "github.com/razshare/frizzante"

func api(withPattern f.WithApiPattern, withHandler f.WithApiHandler) {
	withPattern("GET /")
	withHandler(func(request *f.Request, response *f.Response) {
		// Handle.
		f.SendEcho(response, "Ok.")
	})
}
