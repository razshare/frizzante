package sessions

type Session struct {
	Error string `json:"error"`
	Todos []Todo `json:"todos"`
}

type Todo struct {
	Description string `json:"description"`
	Checked     bool   `json:"checked"`
}
