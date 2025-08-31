package action

import (
	"github.com/razshare/frizzante/internal/cli/generate"
	"github.com/razshare/frizzante/internal/tui/input"
)

func CreateProject(options CreateProjectOptions) (err error) {
	if options.Name == "" {
		options.Name, err = input.Send("give the project a name")
		if err != nil {
			return err
		}
	}

	return generate.Project(generate.ProjectOptions{
		Name: options.Name,
		Efs:  options.Efs,
	})
}
