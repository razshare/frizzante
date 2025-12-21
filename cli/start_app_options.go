package cli

import "github.com/razshare/frizzante/cli/apps"

type StartAppOptions struct {
	App   apps.App
	Query string
}
