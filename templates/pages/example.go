package pages

import f "github.com/razshare/frizzante"

func page(
	withPath func(path string),
	withDocument func(document *f.Document),
	withBaseHandler func(base func(req *f.Request, res *f.Response, doc *f.Document)),
	withActionHandler func(action func(req *f.Request, res *f.Response, doc *f.Document)),
) {
	withPath("/path")
	withDocument(f.DocumentCreate("pageName"))
	withBaseHandler(func(req *f.Request, res *f.Response, doc *f.Document) {
		// Show page.
	})
	withActionHandler(func(req *f.Request, res *f.Response, doc *f.Document) {
		// Modify state.
	})
}
