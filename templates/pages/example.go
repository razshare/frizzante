package pages

import (
	f "github.com/razshare/frizzante"
	//config.Server
)

var Server *f.Server
var guards []f.Guard

func init() {
	Server.LoadController(func(controller *f.Controller) {
		controller.
			WithBase(guards, base).
			WithAction(guards, action)
	})
}

func base(_ *f.Request, res *f.Response) {
	res.SendView(f.NewView(f.RenderModeFull))
}

func action(_ *f.Request, res *f.Response) {
	res.SendView(f.NewView(f.RenderModeFull))
}
