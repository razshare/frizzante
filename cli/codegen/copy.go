package codegen

import (
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
	"regexp"
)

func Copy(cops []CopyInstruction) error {
	for _, cop := range cops {
		if embeds.IsDirectory(cop.Efs, cop.From) {
			ds, err := cop.Efs.ReadDir(cop.From)
			if err != nil {
				return err
			}

			gsloc := make([]CopyInstruction, 0)

			for _, d := range ds {
				gsloc = append(gsloc, CopyInstruction{
					From:      filepath.Join(cop.From, d.Name()),
					To:        filepath.Join(cop.To, d.Name()),
					Overwrite: cop.Overwrite,
				})
			}
			err = Copy(gsloc)
			if err != nil {
				return err
			}
			continue
		}

		from := cop.From
		to := cop.To

		if cop.Overwrite != nil && files.IsFile(to) {
			overwrite, err := cop.Overwrite(to)
			if err != nil {
				return err
			}
			if !overwrite {
				continue
			}
		}

		dat, err := cop.Efs.ReadFile(from)
		if err != nil {
			return err
		}

		cont, err := Parse(string(dat), func(s Section) error {
			for _, m := range s.Mods {
				reg, cerr := regexp.Compile("(?i)" + m.Pattern)
				if cerr != nil {
					return cerr
				}
				*s.Line = reg.ReplaceAllLiteralString(*s.Line, m.Replacement)
			}
			return nil
		})

		if err != nil {
			return err
		}

		dir := filepath.Dir(to)
		if !embeds.IsDirectory(cop.Efs, dir) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
		err = os.WriteFile(to, []byte(cont), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
