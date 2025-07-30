package views

import "rogchap.com/v8go"

type RenderMode int

const (
	RenderModeFull     RenderMode = 0 // Renders on both the server and the client.
	RenderModeServer   RenderMode = 1 // Renders only on the server.
	RenderModeClient   RenderMode = 2 // Renders only on the client.
	RenderModeHeadless RenderMode = 3 // Renders only on the server and omits the base template.
)

type Configuration struct {
}

type View struct {
	Name       string
	Data       map[string]any
	RenderMode RenderMode
	Functions  map[string]v8go.FunctionCallback
	AppRoot    string
	ServerJs   string
	IndexHtml  string
}
