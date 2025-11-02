package select_npm_packages

import "github.com/razshare/frizzante/cli/npm"

type SearchResultMsg struct {
	Packages []npm.PackageInfo
	Error    error
}
