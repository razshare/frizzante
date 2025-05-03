package guards

import f "github.com/razshare/frizzante"

func guard(withHandler func(func(*f.Request, *f.Response, func()))) {
	// Guard.
	withHandler(func(request *f.Request, response *f.Response, pass func()) {
		pass()
	})
}
