package pages

import f "github.com/razshare/frizzante"

func page(
	withPath f.ProvidePagePath,
	withView f.ProvidePageView,
	withBaseHandler f.ProvidePageBaseHandler,
	withActionHandler f.ProvidePageActionHandler,
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
