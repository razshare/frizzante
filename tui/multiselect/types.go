package multiselect

type Model struct {
	Choices  []string
	Cursor   int
	Selected map[int]bool
	Prompt   string
}
