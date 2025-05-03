package pages

import f "github.com/razshare/frizzante"

func page(
	withPath f.WithPagePath,
	withView f.WithPageView,
	withBaseHandler f.WithPageBaseHandler,
	withActionHandler f.WithPageActionHandler,
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
