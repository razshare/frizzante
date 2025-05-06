package api

import f "github.com/razshare/frizzante"

func api(context f.ApiContext) {
	// Context.
	withPattern, withHandler := context()

	// Configure.
	withPattern("GET /")
	withHandler(func(request *f.Request, response *f.Response) {
		// Handle.
		f.ResponseSendMessage(response, "Ok.")
	})
}
