package view

func Data(v View) map[string]any {
	return map[string]any{
		"name":   v.Name,
		"render": v.Render,
		"align":  v.Align,
		"props":  v.Props,
	}
}
