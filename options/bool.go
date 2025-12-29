package options

type Bool struct {
	Name      string
	Shorthand string
	Value     bool
	Usage     string
	Reference *bool
}
