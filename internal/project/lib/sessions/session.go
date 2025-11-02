package sessions

type Session struct {
	Error string `json:"error"`
	Todos []Todo `json:"todos"`
}
