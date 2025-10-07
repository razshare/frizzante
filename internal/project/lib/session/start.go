package session

var Sessions = map[string]*Session{}

func Start(id string) *Session {
	if state, ok := Sessions[id]; ok {
		return state
	}

	Sessions[id] = New()
	return Sessions[id]
}
