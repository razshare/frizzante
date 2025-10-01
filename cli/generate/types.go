package generate

import (
	"embed"

	"github.com/razshare/frizzante/platform"
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
	Platform platform.Platform
	Auto     bool
}

type BunOptions struct {
	Bun      string
	Platform platform.Platform
	Auto     bool
}

type SqlcOptions struct {
	Sqlc     string
	Platform platform.Platform
	Auto     bool
}

type CoreOptions struct {
	App  string
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
	Platform platform.Platform
	Auto     bool
}

type QueriesOptions struct {
	Sqlc     string
	SqlcYaml string
	Platform platform.Platform
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
	App  string
	Efs  embed.FS
	Auto bool
}

type IconsOptions struct {
	Bun  string
	App  string
	Efs  embed.FS
	Auto bool
}

type LinksOptions struct {
	App  string
	Efs  embed.FS
	Auto bool
}

type SessionOptions struct {
	Type string
	Efs  embed.FS
	Auto bool
}
