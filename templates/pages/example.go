package pages

import f "github.com/razshare/frizzante"

func page(
	withPath f.PagePathProvider,
	withView f.PageViewProvider,
	withBaseHandler f.PageBaseHandlerProvider,
	withActionHandler f.PageActionHandlerProvider,
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
