package pages

import f "github.com/razshare/frizzante"

type Controller struct{}

func (_ Controller) Configure(meta func() f.Metadata) f.PageConfiguration {
	return f.PageConfiguration{
		Metadata: meta(),
	}
}

func (_ Controller) Base(req *f.Request, res *f.Response) {
	res.SendView(f.NewView(f.RenderModeFull))
}

func (_ Controller) Action(req *f.Request, res *f.Response) {
	res.SendView(f.NewView(f.RenderModeFull))
}
