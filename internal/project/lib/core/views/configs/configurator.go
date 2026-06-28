package configs

import (
	"embed"
	"log"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

type Configurator func(
	http.ResponseWriter,
	*http.Request,
	views.View,
) (
	http.ResponseWriter,
	*http.Request,
	renders.Render,
	embed.FS,
	*log.Logger,
	*log.Logger,
	views.View,
)
