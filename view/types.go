package view

type RenderMode int

const (
	RenderFull     RenderMode = 0 // RenderFull renders on both the server and the client.
	RenderServer   RenderMode = 1 // RenderServer renders only on the server.
	RenderClient   RenderMode = 2 // RenderClient renders only on the client.
	RenderHeadless RenderMode = 3 // RenderHeadless renders only on the server and omits the base template.
)

type AlignMode int

const (
	AlignReset AlignMode = 0 // AlignReset resets the client view props before injecting given props.
	AlignMerge AlignMode = 1 // AlignMerge merges given props with existing props on the client view.
)

type View struct {
	Name   string
	Title  string
	Render RenderMode
	Align  AlignMode
	Props  map[string]any
}

type Render func(v View) (html string, err error)
