package pages

import f "github.com/razshare/frizzante"

func page(
	withPath func(path string),
	withView func(view *f.View),
	withBaseHandler func(baseHandler func(
		request *f.Request,
		response *f.Response,
		view *f.View,
	)),
	withActionHandler func(actionHandler func(
		request *f.Request,
		response *f.Response,
		view *f.View,
	)),
) {
	withPath("/path")
	withView(f.ViewReference("viewName"))
	withBaseHandler(func(
		request *f.Request,
		response *f.Response,
		view *f.View,
	) {
		// Show page.
	})
	withActionHandler(func(
		request *f.Request,
		response *f.Response,
		view *f.View,
	) {
		// Modify state.
	})
}
