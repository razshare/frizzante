package guards

import f "github.com/razshare/frizzante"

func guard(
	withHandler func(handler func(req *f.Request, res *f.Response, pass func())),
) {
	withHandler(func(req *f.Request, res *f.Response, pass func()) {
		// Guard.
		pass()
	})
}
