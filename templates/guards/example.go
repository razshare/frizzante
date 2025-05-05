package guards

import f "github.com/razshare/frizzante"

func guard(context f.GuardContext) {
	// Context.
	handler := context()

	// Configure.
	handler(func(request *f.Request, response *f.Response, pass func()) {
		// Guard.
		pass()
	})
}
