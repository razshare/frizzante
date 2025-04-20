package indexes

import f "github.com/razshare/frizzante"

func baseHandler(req *f.Request, res *f.Response, page *f.Page) {
	// Show page.
}

func actionHandler(req *f.Request, res *f.Response, page *f.Page) {
	// Modify state.
}

func index(
	withPage func(page string),
	withPath func(path string),
	withBaseHandler func(base func(req *f.Request, res *f.Response, page *f.Page)),
	withActionHandler func(action func(req *f.Request, res *f.Response, page *f.Page)),
) {
	withPage("page")
	withPath("/path")
	withBaseHandler(baseHandler)
	withActionHandler(actionHandler)
}
