package views

type Data struct {
	Name   string     `json:"name"`
	Render RenderMode `json:"render"`
	Align  AlignMode  `json:"align"`
	Props  any        `json:"props"`
}
