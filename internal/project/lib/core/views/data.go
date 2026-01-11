package views

type Data struct {
	IsSnapshot bool       `json:"isSnapshot"`
	Name       string     `json:"name"`
	Render     RenderMode `json:"render"`
	Align      AlignMode  `json:"align"`
	Props      any        `json:"props"`
}
