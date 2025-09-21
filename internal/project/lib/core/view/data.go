package view

func Data(view View) map[string]any {
	return map[string]any{
		"name":   view.Name,
		"render": view.RenderMode,
		"align":  view.AlignMode,
		"props":  view.Props,
	}
}
