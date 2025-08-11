package codegen

import (
	"errors"
	"fmt"
	"github.com/razshare/frizzante/globals"
	"strings"
)

func Parse(c string, b Build) (string, error) {
	var find strings.Builder
	var repl strings.Builder

	ex := func(l string) (Mod, error) {
		find.Reset()
		repl.Reset()
		s := Start
		for _, r := range l {
			switch s {
			case Start:
				switch r {
				case '"':
					s = ReadingOriginalString
				}
			case ReadingOriginalString:
				switch r {
				case '\\':
					s = Escaping
				case '"':
					s = DoneReadingOriginalString
				default:
					find.WriteRune(r)
				}
			case DoneReadingOriginalString:
				switch r {
				case ' ':
					s = ExpectingReplacementString
				}
			case ExpectingReplacementString:
				switch r {
				case '"':
					s = ReadingReplacementString
				default:
				}
			case ReadingReplacementString:
				switch r {
				case '\\':
					s = Escaping
				case '"':
					s = DoneReadingReplacementString
					return Mod{
						Pattern:     find.String(),
						Replacement: repl.String(),
					}, nil
				default:
					repl.WriteRune(r)
				}
			case Escaping:
				// Noop.
				continue
			default:
				switch r {
				case ' ':
					// Noop.
					continue
				}
				s = Invalid
				return Mod{}, fmt.Errorf("invalid state, expecting \",\", received \"%c\" instead", r)
			}
		}
		return Mod{}, errors.New("invalid mod")
	}

	ms := make([]Mod, 0)
	msl := 0
	ml := make([]Mod, 0)
	mll := 0
	o := len(globals.CodegenModHint)
	var sb strings.Builder

	for _, l := range strings.Split(c, "\n") {
		trmd := strings.TrimSpace(l)
		if strings.HasPrefix(trmd, globals.CodegenModsHint) {
			mod, err := ex(l[o:])
			if err != nil {
				return "", err
			}
			ms = append(ms, mod)
			msl++
			continue
		} else if strings.HasPrefix(trmd, globals.CodegenModHint) {
			mod, err := ex(l[o:])
			if err != nil {
				return "", err
			}
			ml = append(ml, mod)
			mll++
			continue
		}

		if msl > 0 {
			err := b(Section{
				Mods: ms,
				Line: &l,
			})
			if err != nil {
				return "", err
			}
		}

		if mll > 0 {
			err := b(Section{
				Mods: append(ms, ml...),
				Line: &l,
			})
			if err != nil {
				return "", err
			}
			ml = make([]Mod, 0)
			mll = 0
		}

		sb.WriteString(l + "\n")

	}

	return sb.String(), nil
}
