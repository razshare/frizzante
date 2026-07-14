package routes

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

type Handler = func(scope scopes.Scope, request *http.Request, writer http.ResponseWriter)
