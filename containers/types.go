package containers

import (
	"embed"
	"github.com/dop251/goja"
)

type ViewContainer struct {
	MaximumProgramCounter uint64
	MaximumRuntimeCounter uint64
	ProgramChannel        chan *goja.Program
	RuntimeChannel        chan *goja.Runtime
	Stop                  bool
	AppRoot               string
	ServerJs              string
	IndexHtml             string
	IndexHtmlCache        string
	Efs                   embed.FS
}
