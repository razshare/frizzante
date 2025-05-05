package pages

import f "github.com/razshare/frizzante"

func page(
	withPath f.ConfigurePagePath,
	withView f.ConfigurePageView,
	withBase f.ConfigurePageBase,
	withAction f.ConfigurePageAction,
) {
	withPath("/path")
	withView(f.ViewReference("ViewName"))
	withBase(func(request *f.Request, response *f.Response, view *f.View) {
		// Show page.
	})
	withAction(func(request *f.Request, response *f.Response, view *f.View) {
		// Modify state.
	})
}
