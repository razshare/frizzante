package sessions

type Session struct {
	Todos []Todo
}

type Todo struct {
	Description string
	Checked     bool
}
