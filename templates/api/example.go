package api

import f "github.com/razshare/frizzante"

func api(context f.ApiContext) {
	// Context.
	pattern, handler := context()

	// Configure.
	pattern("GET /")
	handler(func(request *f.Request, response *f.Response) {
		// Handle.
		f.SendEcho(response, "Ok.")
	})
}
