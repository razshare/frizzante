package codegen

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
	"regexp"
)

func Copy(gs []CopyInstruction) error {
	for _, g := range gs {
		if embeds.IsDirectory(cli.Efs, g.From) {
			ds, err := cli.Efs.ReadDir(g.From)
			if err != nil {
				return err
			}

			gsloc := make([]CopyInstruction, 0)

			for _, d := range ds {
				gsloc = append(gsloc, CopyInstruction{
					From:      filepath.Join(g.From, d.Name()),
					To:        filepath.Join(g.To, d.Name()),
					Overwrite: g.Overwrite,
				})
			}
			err = Copy(gsloc)
			if err != nil {
				return err
			}
			continue
		}

		from := g.From
		to := g.To

		if g.Overwrite != nil && files.IsFile(to) {
			overwrite, err := g.Overwrite(to)
			if err != nil {
				return err
			}
			if !overwrite {
				continue
			}
		}

		dat, err := cli.Efs.ReadFile(from)
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
		if !embeds.IsDirectory(cli.Efs, dir) {
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
