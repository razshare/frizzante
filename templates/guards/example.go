package guards

import f "github.com/razshare/frizzante"

func apiHandler(req *f.Request, res *f.Response, pass func()) {
	// Guard api.
	pass()
}

func pageHandler(req *f.Request, res *f.Response, page *f.Page, pass func()) {
	// Guard page.
	pass()
}

func guard(
	withApiHandler func(func(req *f.Request, res *f.Response, pass func())),
	withPageHandler func(func(req *f.Request, res *f.Response, page *f.Page, pass func())),
) {
	withApiHandler(apiHandler)
	withPageHandler(pageHandler)
}
