package guards

import (
	"net/http"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/scopes"
)

type Handler func(scope scopes.Scope, request *http.Request, writer http.ResponseWriter, allow func())
