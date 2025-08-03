package views

import (
	"embed"
	"github.com/dop251/goja"
	"sync"
)

type RenderMode int

const (
	RenderModeFull     RenderMode = 0 // Renders on both the server and the client.
	RenderModeServer   RenderMode = 1 // Renders only on the server.
	RenderModeClient   RenderMode = 2 // Renders only on the client.
	RenderModeHeadless RenderMode = 3 // Renders only on the server and omits the base template.
)

type View struct {
	Name       string
	RenderMode RenderMode
	Data       map[string]any
	Container  *Container
}

type Container struct {
	*ContainerConfiguration
	IndexHtmlCache string
	Runtime        *goja.Runtime
	Render         goja.Callable
	Mutex          *sync.Mutex
}

type ContainerConfiguration struct {
	Efs        embed.FS
	AppRoot    string
	ServerJs   string
	IndexHtml  string
	ReadMode   ViewReadMode
	BundleMode ViewBundleMode
	Mutex      sync.Mutex
}
