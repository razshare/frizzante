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
	Auto     bool
	Platform platform.Platform
}

type BunOptions struct {
	Bun      string
	Auto     bool
	Platform platform.Platform
}

type SqlcOptions struct {
	Sqlc     string
	Auto     bool
	Platform platform.Platform
}

type CoreOptions struct {
	App  string
	Auto bool
	Lib  string
	Efs  embed.FS
}

type DatabaseOptions struct {
	Efs      embed.FS
	Generate string
	Auto     bool
	Go       string
	Sqlc     string
	Lib      string
	Platform platform.Platform
}

type QueriesOptions struct {
	Auto     bool
	Sqlc     string
	Lib      string
	Platform platform.Platform
}

type DownloadOptions struct {
	Url  string
	Auto bool
}

type CopyOptions struct {
	From string
	To   string
	Auto bool
}

type ProjectOptions struct {
	Auto bool
	Name string
	Efs  embed.FS
}

type EmbeddedZipOptions struct {
	Auto     bool
	FileName string
	Efs      embed.FS
}

type FormsOptions struct {
	App  string
	Auto bool
	Lib  string
	Efs  embed.FS
}

type LinksOptions struct {
	App  string
	Auto bool
	Lib  string
	Efs  embed.FS
}

type SessionOptions struct {
	Auto bool
	Lib  string
	Efs  embed.FS
}
