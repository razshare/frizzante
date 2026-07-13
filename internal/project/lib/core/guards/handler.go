package guards

import (
	"net/http"
)

type Handler = func(scope Scope, request *http.Request, writer http.ResponseWriter, allow func())
