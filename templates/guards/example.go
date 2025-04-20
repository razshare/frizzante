package guards

import f "github.com/razshare/frizzante"

func guardHandler(req *f.Request, res *f.Response, pass func()) {
	// Guard.
	pass()
}

func guard(
	withGuardHandler func(
		guardHandler func(
			req *f.Request,
			res *f.Response,
			pass func(),
		),
	),
) {
	withGuardHandler(guardHandler)
}
