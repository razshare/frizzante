package session

var Sessions = map[string]*Session{}

func Start(id string) *Session {
	v, ok := Sessions[id]
	if !ok {
		Sessions[id] = &Session{Todos: []Todo{
			{Checked: false, Description: "Pet the cat."},
			{Checked: false, Description: "Do laundry"},
			{Checked: false, Description: "Pet the cat."},
			{Checked: false, Description: "Cook"},
			{Checked: false, Description: "Pet the cat."},
		}}
		return Sessions[id]
	}
	return v
}
