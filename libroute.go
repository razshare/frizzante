package frizzante

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Route struct {
	server  *Server
	view    string
	handler func(request *Request, response *Response)
	guards  []func(request *Request, response *Response, pass func())
}

// routeCreate creates a route configuration from a callback function.
func routeCreate(
	handler func(request *Request, response *Response),
	guards []func(request *Request, response *Response, pass func()),
) *Route {
	return &Route{
		view: "",
		handler: func(request *Request, response *Response) {
			pass := 0 == len(guards)
			for _, guard := range guards {
				guard(request, response, func() { pass = true })

				if !pass {
					break
				}
			}

			if pass {
				handler(request, response)
			}
		},
	}
}

// routeCreateWithView creates a route configuration from a callback function, just like routeCreate.
//
// Unlike routeCreate, routeCreateWithView also creates a View, which is used to automatically
// to serve a svelte view after invoking callback.
//
// Generally speaking, you should never manually invoke ResponseSendMessage or similar functions.
//
// However, it is safe to invoke receive functions, like RequestReceiveHeader, RequestReceiveCookie, etc.
func routeCreateWithView[T any](
	handler func(request *Request, response *Response, view *View[T]),
	guards []func(request *Request, response *Response, pass func()),
	view *View[T],
) *Route {
	var pattern string
	return &Route{
		view: view.name,
		handler: func(
			request *Request,
			response *Response,
		) {
			viewScoped := &View[T]{
				name:               view.name,
				functions:          view.functions,
				embeddedFileSystem: view.embeddedFileSystem,
				dataInitializer:    view.dataInitializer,
				Render:             view.Render,
				Data:               view.dataInitializer(),
			}

			pass := 0 == len(guards)
			for _, guard := range guards {
				guard(request, response, func() { pass = true })

				if !pass {
					break
				}
			}

			if pass {
				handler(request, response, viewScoped)
			}

			if nil != response.navigate {
				ResponseSendRedirect(response, response.navigate.Location, http.StatusFound)
				return
			}

			if "" != response.header.Get("Location") {
				return
			}

			if nil == viewScoped {
				NotifierSendError(request.server.notifier, fmt.Errorf("svelte page handler `%s` returned a nil page", pattern))
				return
			}

			if RequestVerifyAccept(request, "application/json") {
				data, marshalError := json.Marshal(viewScoped.Data)
				if nil != marshalError {
					NotifierSendError(request.server.notifier, marshalError)
					return
				}
				ResponseSendHeader(response, "Content-Type", "application/json")
				ResponseSendMessage(response, string(data))
				return
			}

			ResponseSendView(response, viewScoped)
		},
	}
}
