package views

type RenderMode int

const (
	RenderModeFull     RenderMode = 0 // RenderModeFull renders on both the server and the client.
	RenderModeServer   RenderMode = 1 // RenderModeServer renders only on the server.
	RenderModeClient   RenderMode = 2 // RenderModeClient renders only on the client.
	RenderModeHeadless RenderMode = 3 // ModeHeadless renders only on the server and omits the base template.
)
