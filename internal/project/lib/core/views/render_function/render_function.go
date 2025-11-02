package render_function

import "github.com/razshare/frizzante/internal/project/lib/core/views"

type RenderFunction = func(view views.View) (head string, body string, err error)
