package generate

import "github.com/razshare/frizzante/platforms"

type QueriesOptions struct {
	Sqlc     string
	SqlcYaml string
	Platform platforms.Platform
	Auto     bool
}
