package pages

import f "github.com/razshare/frizzante"

type pageData struct {
}

type pageController struct {
	f.PageController
}

func (_ pageController) Configure() f.PageConfiguration {
	return f.PageConfiguration{
		Path: "/path",
	}
}

func (_ pageController) Base(req *f.Request, res *f.Response) {
	res.SendView(f.NewViewWithData(f.RenderModeFull, pageData{}))
}

func (_ pageController) Action(req *f.Request, res *f.Response) {
	res.SendView(f.NewViewWithData(f.RenderModeFull, pageData{}))
}
