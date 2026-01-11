//go:build snapshot_servers

package views

func NewData(view View) Data {
	return Data{
		IsSnapshot: true,
		Name:       view.Name,
		Render:     view.RenderMode,
		Align:      view.AlignMode,
		Props:      view.Props,
	}
}
