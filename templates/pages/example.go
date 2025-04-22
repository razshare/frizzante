package pages

import f "github.com/razshare/frizzante"

func page(
	withPath func(path string),
	withDocument func(document *f.Document),
	withBaseHandler func(baseHandler func(request *f.Request, response *f.Response, document *f.Document)),
	withActionHandler func(actionHandler func(request *f.Request, response *f.Response, document *f.Document)),
) {
	withPath("/path")
	withDocument(f.DocumentCreate("pageName"))
	withBaseHandler(func(request *f.Request, response *f.Response, document *f.Document) {
		// Show page.
	})
	withActionHandler(func(request *f.Request, response *f.Response, document *f.Document) {
		// Modify state.
	})
}
