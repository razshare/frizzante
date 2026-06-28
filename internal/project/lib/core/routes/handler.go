package routes

import (
	"net/http"
)

type Handler = func(request *http.Request, writer http.ResponseWriter)
