package javascript

type RenderFunction = func(options RenderFunctionOptions) (head string, body string, err error)
