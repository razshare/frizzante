package pages

import f "github.com/razshare/frizzante"

func page(page *f.Page) {
	f.PageWithPath(page, "/path")
	f.PageWithView(page, f.ViewReference("ViewName"))
	f.PageWithBaseHandler(page, func(request *f.Request, response *f.Response, view *f.View) {
		// Show page.
	})
	f.PageWithActionHandler(page, func(request *f.Request, response *f.Response, view *f.View) {
		// Modify state.
	})
}
