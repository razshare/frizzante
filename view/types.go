package view

type RenderMode int

const (
	// RenderModeFull renders on both the server and the client.
	RenderModeFull RenderMode = 0
	// RenderModeServer renders only on the server.
	RenderModeServer RenderMode = 1
	// RenderModeClient renders only on the client.
	RenderModeClient RenderMode = 2
	// RenderModeHeadless renders only on the server and omits the base template.
	RenderModeHeadless RenderMode = 3
)

type View struct {
	Name       string
	RenderMode RenderMode
	Data       map[string]any
}
