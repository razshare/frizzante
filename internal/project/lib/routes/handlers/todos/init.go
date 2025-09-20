//go:build types

package todos

import "github.com/razshare/frizzante/internal/project/lib/core/types"

func init() { types.Generate[Props]() }
