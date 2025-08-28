//gen:mod "memory" "session"
package memory

type State struct {
	Todos []Todo
}

type Todo struct {
	Checked     bool
	Description string
}
