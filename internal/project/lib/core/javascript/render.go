package javascript

type Render = func(options RenderOptions) (head string, body string, err error)
