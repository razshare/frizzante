package containers

import (
	"embed"
	"github.com/dop251/goja"
)

type ViewContainer struct {
	Stop                  bool
	MaximumProgramCounter uint64
	MaximumRuntimeCounter uint64
	AppRoot               string
	ServerJs              string
	IndexHtml             string
	IndexHtmlCache        string
	ProgramChannel        chan *goja.Program
	RuntimeChannel        chan *goja.Runtime
	Efs                   embed.FS
}
