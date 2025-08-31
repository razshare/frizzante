package session

type State struct {
	Todos []Todo
}

type Todo struct {
	Checked     bool
	Description string
}
