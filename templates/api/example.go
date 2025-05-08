package api

import f "github.com/razshare/frizzante"

func api(api *f.Api) {
	f.ApiWithPattern(api, "GET /")
	f.ApiWithRequestHandler(api, func(request *f.Request, response *f.Response) {
		// Handle.
		f.ResponseSendMessage(response, "Ok.")
	})
}
