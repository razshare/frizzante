package renders

type Render = func(options RenderOptions) (html string, err error)
