package files

import (
	"github.com/razshare/frizzante/codegen"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func Generate(gs []codegen.Generation) error {
	for _, g := range gs {
		if IsDirectory(g.From) {
			ds, err := os.ReadDir(g.From)
			if err != nil {
				return err
			}

			gsloc := make([]codegen.Generation, 0)

			for _, d := range ds {
				gsloc = append(gsloc, codegen.Generation{
					From: filepath.Join(g.From, d.Name()),
					To:   filepath.Join(g.To, d.Name()),
				})
			}
			err = Generate(gsloc)
			if err != nil {
				return err
			}
			continue
		}

		from := g.From
		to := g.To

		dat, err := os.ReadFile(from)
		if err != nil {
			log.Fatal(err)
		}

		cont, err := codegen.Parse(string(dat), func(s codegen.Section) error {
			for _, m := range s.Mods {
				reg, cerr := regexp.Compile(m.Pattern)
				if cerr != nil {
					return cerr
				}
				*s.Line = reg.ReplaceAllString(*s.Line, m.Replacement)
			}
			return nil
		})

		if err != nil {
			return err
		}

		dir := filepath.Dir(to)
		if !IsDirectory(dir) {
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
