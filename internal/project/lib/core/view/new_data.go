package view_

func NewData(view View) Data {
	return Data{
		Name:   view.Name,
		Render: view.RenderMode,
		Align:  view.AlignMode,
		Props:  view.Props,
	}
}
