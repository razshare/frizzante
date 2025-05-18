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

func (_ pageController) Base(request *f.Request, response *f.Response) {
	response.SendView(f.NewView(pageData{}))
}

func (_ pageController) Action(request *f.Request, response *f.Response) {
	response.SendView(f.NewView(pageData{}))
}
