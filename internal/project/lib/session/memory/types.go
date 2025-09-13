package session

type Session struct {
	Todos []Todo
}

type Todo struct {
	Checked     bool
	Description string
}
