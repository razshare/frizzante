package guards

import f "github.com/razshare/frizzante"

func guard(guard *f.Guard) {
	f.GuardWithHandler(guard, func(request *f.Request, response *f.Response, pass func()) {
		// Guard.
		pass()
	})
}
