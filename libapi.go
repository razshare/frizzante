package frizzante

type Api struct {
	patterns []string
	handler  func(request *Request, response *Response)
	guards   []func(request *Request, response *Response, pass func())
}

// ApiBuilder builds an api.
type ApiBuilder = func(api *Api)

func ApiCreate() *Api {
	return &Api{
		patterns: []string{},
		guards:   []func(request *Request, response *Response, pass func()){},
		handler: func(request *Request, response *Response) {
			// Noop.
		},
	}
}

// ApiWithPattern adds a pattern.
func ApiWithPattern(self *Api, pattern string) {
	self.patterns = append(self.patterns, pattern)
}

// ApiWithRequestHandler sets the request handler.
func ApiWithRequestHandler(self *Api, handler func(request *Request, response *Response)) {
	self.handler = handler
}

// ApiWithGuardHandler adds a guard handler.
func ApiWithGuardHandler(self *Api, handler func(request *Request, response *Response, pass func())) {
	self.guards = append(self.guards, handler)
}
