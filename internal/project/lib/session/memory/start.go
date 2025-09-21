package session

var States = map[string]*State{}

func Start(id string) *State {
	v, ok := States[id]
	if !ok {
		States[id] = New()
		return States[id]
	}
	return v
}
