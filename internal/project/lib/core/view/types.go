package view

type RenderMode int

const (
	RenderModeFull     RenderMode = 0 // RenderModeFull renders on both the server and the client.
	RenderModeServer   RenderMode = 1 // RenderModeServer renders only on the server.
	RenderModeClient   RenderMode = 2 // RenderModeClient renders only on the client.
	RenderModeHeadless RenderMode = 3 // ModeHeadless renders only on the server and omits the base template.
)

type AlignMode int

const (
	AlignModeReset AlignMode = 0 // AlignModeReset resets the client view props before injecting given props.
	AlignModeMerge AlignMode = 1 // AlignModeMerge merges given props with existing props on the client view.
)

type View struct {
	Name       string
	Title      string
	Props      map[string]any
	AlignMode  AlignMode
	RenderMode RenderMode
}
