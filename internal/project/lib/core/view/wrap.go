package view

func Wrap(view View) Data {
	return Data{
		Name:   view.Name,
		Render: view.RenderMode,
		Align:  view.AlignMode,
		Props:  view.Props,
	}
}
