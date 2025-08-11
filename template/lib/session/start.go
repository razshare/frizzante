package session

//gen:mods "session" "Session"
//gen:mods "start" "Start"
var sessions = map[string]*session{}

func start(id string) *session {
	v, ok := sessions[id]
	if !ok {
		sessions[id] = &session{}
		return sessions[id]
	}
	return v
}
