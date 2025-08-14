//gen:mod "memory" "session"
package memory

//gen:mod "state" "State"
type state struct {
	//gen:mod "todo" "Todo"
	Todos []todo
}

//gen:mod "todo" "Todo"
type todo struct {
	Checked     bool
	Description string
}
