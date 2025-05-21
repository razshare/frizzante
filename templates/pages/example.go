package pages

import f "github.com/razshare/frizzante"

var pageName = f.
	NewPageController().
	WithBase(pageBase).
	WithAction(pageAction)

func pageBase(_ *f.Request, res *f.Response) {
	res.SendView(f.NewView(f.RenderModeFull))
}

func pageAction(_ *f.Request, res *f.Response) {
	res.SendView(f.NewView(f.RenderModeFull))
}
