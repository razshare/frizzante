package guards

import f "github.com/razshare/frizzante"

func guard(withHandler f.GuardHandlerProvider) {
	withHandler(func(request *f.Request, response *f.Response, pass func()) {
		// Guard.
		pass()
	})
}
