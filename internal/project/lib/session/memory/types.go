package session

type State struct {
	Todos []Todo `json:"todos"`
}

type Todo struct {
	Description string `json:"description"`
	Checked     bool   `json:"checked"`
}
