package guards

import (
	"net/http"
)

type Handler = func(request *http.Request, writer http.ResponseWriter, allow func())
