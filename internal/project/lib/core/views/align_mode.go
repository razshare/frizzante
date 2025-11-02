package views

type AlignMode int

const (
	AlignModeReset AlignMode = 0 // AlignModeReset resets the client view props before injecting given props.
	AlignModeMerge AlignMode = 1 // AlignModeMerge merges given props with existing props on the client view.
)
