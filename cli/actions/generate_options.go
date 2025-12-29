package actions

import "embed"

type GenerateOptions struct {
	Generation   string
	Go           string
	Air          string
	Bun          string
	Sqlc         string
	SqlcYaml     string
	Database     string
	DatabaseType string
	Tags         string
	Efs          embed.FS
}
