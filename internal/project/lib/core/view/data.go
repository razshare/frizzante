package view

func Data(v View) map[string]any {
	return map[string]any{
		"name":   v.Name,
		"render": v.RenderMode,
		"align":  v.AlignMode,
		"props":  v.Props,
	}
}
