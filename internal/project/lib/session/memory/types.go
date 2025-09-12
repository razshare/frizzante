package session

type Session struct {
	Mode  Mode
	Todos []Todo
}

type Todo struct {
	Checked     bool
	Description string
}

type Mode uint

const ModeToggle Mode = 0
const ModeRemove Mode = 1
const ModeAdd Mode = 2
