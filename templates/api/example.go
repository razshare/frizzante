package api

import f "github.com/razshare/frizzante"

func api(withPattern f.ApiPatternProvider, withHandler f.ApiHandlerProvider) {
	withPattern("GET /")
	withHandler(func(request *f.Request, response *f.Response) {
		// Handle.
		f.SendEcho(response, "Ok.")
	})
}
