package generate

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
	"github.com/razshare/frizzante/tui/spinners"
)

func Queries(options QueriesOptions) (err error) {
	yamlFileName := options.SqlcYaml

	if yamlFileName != "" && !files.IsFile(yamlFileName) {
		messages.Infof("%s not found", yamlFileName)
	}

	if yamlFileName == "" {
		var items []string
		if items, err = files.ReadDirectory("lib"); err != nil {
			return
		}

		names := make([]string, 0)
		for _, item := range items {
			if strings.HasSuffix(item, string(filepath.Separator)+"sqlc.yaml") {
				names = append(names, item)
			}
		}

		choices := make([]search.Choice, len(names))
		for index, name := range names {
			choices[index] = search.Choice{Id: name}
		}

		choices = append(choices, search.Choice{Id: "other", Description: "other"})

		if options.Auto {
			yamlFileName = choices[0].Id
		} else {
			yamlFileName, err = select_one.Sendf(choices, "where is your sqlc.yaml file located?")
		}
	}

	baseDirectory := filepath.Dir(yamlFileName)

	if _, err = exec.LookPath(options.Sqlc); err != nil && !files.IsFile(options.Sqlc) {
		if err = Sqlc(SqlcOptions{
			Sqlc:     options.Sqlc,
			Platform: options.Platform,
			Auto:     options.Auto,
		}); err != nil {
			return
		}
	}

	if !files.IsFile(yamlFileName) {
		return fmt.Errorf("%s not found", yamlFileName)
	}

	var sqlc string
	if files.IsFile(options.Sqlc) {
		if sqlc, err = filepath.Rel(baseDirectory, options.Sqlc); err != nil {
			return err
		}
	} else if sqlc, err = exec.LookPath(options.Sqlc); err != nil {
		sqlc = options.Sqlc
	}

	spin := spinners.New("generating queries")

	go spinners.Start(spin)
	if !messages.Command(baseDirectory, os.Environ(), sqlc, "generate") {
		spinners.Stop(spin)
		err = errors.New("could not generate queries")
		return
	}
	spinners.Stop(spin)

	if err = FixImports(FixImportsOptions{Directory: baseDirectory}); err != nil {
		return
	}

	messages.Success(filepath.Join(baseDirectory, "queries.go"))

	return
}
