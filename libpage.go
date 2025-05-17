package frizzante

import (
	"fmt"
)

type Page[T any] struct {
	server *Server
	path   string
	name   string
	view   *View[T]
	base   func(request *Request, response *Response, view *View[T])
	action func(request *Request, response *Response, view *View[T])
	guards []func(request *Request, response *Response, pass func())
}

// PageBuilder builds a page.
type PageBuilder[T any] = func(page *Page[T])

type PageMetadata struct {
	Path     string `json:"path"`
	ViewName string `json:"viewName"`
}

var pages = map[string]PageMetadata{}

// PageCreate creates a page.
func PageCreate[T any]() *Page[T] {
	page := &Page[T]{
		guards: []func(request *Request, response *Response, pass func()){},
		base: func(request *Request, response *Response, view *View[T]) {
			// Noop.
		},
		action: func(request *Request, response *Response, view *View[T]) {
			// Noop.
		},
	}
	return page
}

// PageWithView sets the view by name.
func PageWithView[T any](self *Page[T], name string, initializer ViewDataInitializer[T]) {
	view, viewError := ViewCreate[T](name, initializer)
	if nil != viewError {
		NotifierSendError(self.server.notifier, viewError)
		return
	}

	ViewWithEmbeddedFileSystem(view, self.server.embeddedFileSystem)

	self.view = view
}

// PageWithServer sets the server.
func PageWithServer[T any](self *Page[T], server *Server) {
	self.server = server
}

// PageWithPath sets the path.
func PageWithPath[T any](self *Page[T], path string) {
	if "" == path {
		NotifierSendError(self.server.notifier, fmt.Errorf("page path cannot be empty"))
		return
	}
	self.path = path
}

// PageWithName gives the page a name, allowing the client to navigate to it.
func PageWithName[T any](self *Page[T], name string) {
	if !KeyIsSafe(name) {
		NotifierSendError(self.server.notifier, fmt.Errorf("page name `%s` is not safe", name))
		return
	}
	self.name = name
}

// PageWithBaseHandler sets the base.
func PageWithBaseHandler[T any](self *Page[T], handler func(request *Request, response *Response, view *View[T])) {
	self.base = handler
}

// PageWithActionHandler sets the action handler.
func PageWithActionHandler[T any](self *Page[T], handler func(request *Request, response *Response, view *View[T])) {
	self.action = handler
}

// PageWithGuardHandler adds a guard handler.
func PageWithGuardHandler[T any](self *Page[T], handler func(request *Request, response *Response, pass func())) {
	self.guards = append(self.guards, handler)
}
