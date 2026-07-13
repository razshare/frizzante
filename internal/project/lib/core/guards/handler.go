package guards

import (
	"net/http"
)

type Handler = func(id uint64, request *http.Request, writer http.ResponseWriter, allow func())
