//gen:mod "memory" "session"
package memory

//gen:mod "state" "State"
//gen:mod "sessions" "Sessions"
var sessions = map[string]*state{}

//gen:mod "start" "Start"
//gen:mod "state" "State"
func start(id string) *state {
	//gen:mod "sessions" "Sessions"
	v, ok := sessions[id]
	if !ok {
		//gen:mod "newState" "New"
		//gen:mod "sessions" "Sessions"
		sessions[id] = newState()
		//gen:mod "sessions" "Sessions"
		return sessions[id]
	}
	return v
}
