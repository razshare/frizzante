package generate

import (
	"database/sql"
	"embed"

	"github.com/razshare/frizzante/platforms"
)

type State uint64

const Start State = 0
const EscapingOriginal State = 98
const EscapingReplacement State = 99
const ReadingOriginalString State = 100
const DoneReadingOriginalString State = 101
const ReadingReplacementString State = 200
const DoneReadingReplacementString State = 201
const ExpectingReplacementString State = 900
const Invalid State = 1000

type Mod struct {
	Pattern     string
	Replacement string
}

type Submit func(char string)
type Build func(block Block) error
type Block struct {
	Mods []Mod
	Line *string
}

type Install func(to string) (bool, error)
type Evict func() error

type AirOptions struct {
	Air      string
	Platform platforms.Platform
	Auto     bool
}

type AirConfigOptions struct {
	Tags []string
	Efs  embed.FS
}

type BunOptions struct {
	Bun      string
	Platform platforms.Platform
	Auto     bool
}

type SqlcOptions struct {
	Sqlc     string
	Platform platforms.Platform
	Auto     bool
}

type CoreOptions struct {
	Efs  embed.FS
	Auto bool
}

type SecurityOptions struct {
	Efs  embed.FS
	Auto bool
}

type DatabaseOptions struct {
	Generate string
	Go       string
	Type     string
	Sqlc     string
	Efs      embed.FS
	Platform platforms.Platform
	Auto     bool
}

type SchemaOptions struct {
	Sqlc     string
	SqlcYaml string
	Database *sql.DB
	Platform platforms.Platform
	Auto     bool
}

type QueriesOptions struct {
	Sqlc     string
	SqlcYaml string
	Platform platforms.Platform
	Auto     bool
}

type DefinitionsOptions struct {
	Go string
}

type MigrationOptions struct {
	Sqlc     string
	SqlcYaml string
	Platform platforms.Platform
	Auto     bool
}

type DownloadOptions struct {
	Url  string
	Auto bool
}

type CopyOptions struct {
	Ignore []string
	From   string
	To     string
	Efs    embed.FS
	Auto   bool
}

type ProjectOptions struct {
	Name string
	Go   string
	Efs  embed.FS
	Auto bool
}

type FixImportsOptions struct {
	Directory string
}

type EmbeddedZipOptions struct {
	FileName string
	Efs      embed.FS
	Auto     bool
}

type FormsOptions struct {
	Efs  embed.FS
	Auto bool
}

type IconsOptions struct {
	Bun string

	Efs  embed.FS
	Auto bool
}

type LinksOptions struct {
	Efs  embed.FS
	Auto bool
}

type SessionsOptions struct {
	Type string
	Efs  embed.FS
	Auto bool
}
