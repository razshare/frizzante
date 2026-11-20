package databases

import (
	"embed"
)

type ListMigrationsOptions struct {
	Efs       embed.FS
	Direction MigrationDirection
	Exists    func(version string) (exists bool, err error)
}
