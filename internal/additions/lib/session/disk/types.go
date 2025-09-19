package session

type State struct {
	Todos []Todo
}

type Todo struct {
	Description string
	Checked     bool
}
