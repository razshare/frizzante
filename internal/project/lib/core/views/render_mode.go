package views

type RenderMode int

const (
	RenderModeFull     RenderMode = 0 // RenderModeFull renders on both the server and the http.
	RenderModeServer   RenderMode = 1 // RenderModeServer renders only on the server.
	RenderModeClient   RenderMode = 2 // RenderModeClient renders only on the http.
	RenderModeHeadless RenderMode = 3 // ModeHeadless renders only on the server and omits the base template.
)
