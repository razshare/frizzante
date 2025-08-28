//gen:mod "memory" "session"
package memory

var Sessions = map[string]*State{}

func Start(id string) *State {
	v, ok := Sessions[id]
	if !ok {
		Sessions[id] = New()
		return Sessions[id]
	}
	return v
}
