package pages

import f "github.com/razshare/frizzante"

func page(context f.PageContext) {
	// Context.
	path, view, base, action := context()

	// Configure.
	path("/path")
	view(f.ViewReference("ViewName"))
	base(func(request *f.Request, response *f.Response, view *f.View) {
		// Show page.
	})
	action(func(request *f.Request, response *f.Response, view *f.View) {
		// Modify state.
	})
}
