package frizzante

type Guard = func(req *Request, res *Response) bool

func AllGuardsPass(req *Request, res *Response, guards ...Guard) bool {
	for _, guard := range guards {
		if !guard(req, res) {
			return false
		}
	}
	return true
}
