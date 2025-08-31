package session

func New() *State {
	return &State{Todos: []Todo{
		{Checked: false, Description: "Pet the cat."},
		{Checked: false, Description: "Do laundry"},
		{Checked: false, Description: "Pet the cat."},
		{Checked: false, Description: "Cook"},
		{Checked: false, Description: "Pet the cat."},
	}}
}
