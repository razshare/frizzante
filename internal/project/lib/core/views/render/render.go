package render

type Render = func(options Options) (html string, err error)
