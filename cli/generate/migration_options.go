package generate

import "github.com/razshare/frizzante/platforms"

type MigrationOptions struct {
	Sqlc     string
	SqlcYaml string
	Platform platforms.Platform
	Auto     bool
}
