package routes

import (
	"net/http"
)

type Handler = func(scope Scope, request *http.Request, writer http.ResponseWriter)
