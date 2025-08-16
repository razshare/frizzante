//gen:mod "disk" "session"
package disk

//gen:mod "newState" "New"
//gen:mod "state" "State"
func newState() *state {
	//gen:mod "state" "State"
	//gen:mod "todo" "Todo"
	return &state{Todos: []todo{
		{Checked: false, Description: "Pet the cat."},
		{Checked: false, Description: "Do laundry"},
		{Checked: false, Description: "Pet the cat."},
		{Checked: false, Description: "Cook"},
		{Checked: false, Description: "Pet the cat."},
	}}
}
