package pages

import f "github.com/razshare/frizzante"

var pageName = f.
	NewPageController().
	WithBase(base).
	WithAction(action)

type data struct{}

func base(_ *f.Request, res *f.Response) {
	res.SendView(f.NewViewWithData(f.RenderModeFull, data{}))
}

func action(_ *f.Request, res *f.Response) {
	res.SendView(f.NewViewWithData(f.RenderModeFull, data{}))
}
