package guards

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

type Handler = func(http *scopes.Http, allow func())
