package pages

import f "github.com/razshare/frizzante"

func page(
	withPath func(string),
	withView func(*f.View),
	withBaseHandler func(func(*f.Request, *f.Response, *f.View)),
	withActionHandler func(func(*f.Request, *f.Response, *f.View)),
) {
	withPath("/path")
	withView(f.ViewReference("ViewName"))
	withBaseHandler(func(request *f.Request, response *f.Response, view *f.View) {
		// Show page.
	})
	withActionHandler(func(request *f.Request, response *f.Response, view *f.View) {
		// Modify state.
	})
}
