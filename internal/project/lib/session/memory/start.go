package session

var States = map[string]*State{}

func Start(id string) *State {
	if state, ok := States[id]; ok {
		return state
	}

	States[id] = New()
	return States[id]
}
