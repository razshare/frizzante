package pages

import f "github.com/razshare/frizzante"

func page(context f.PageContext) {
	// Context.
	withPath, withView, withBase, withAction := context()

	// Configure.
	withPath("/path")
	withView(f.ViewReference("ViewName"))
	withBase(func(request *f.Request, response *f.Response, view *f.View) {
		// Show page.
	})
	withAction(func(request *f.Request, response *f.Response, view *f.View) {
		// Modify state.
	})
}
