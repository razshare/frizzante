package guards

import f "github.com/razshare/frizzante"

func guard(context f.GuardContext) {
	// Context.
	withHandler := context()

	// Configure.
	withHandler(func(request *f.Request, response *f.Response, pass func()) {
		// Guard.
		pass()
	})
}
