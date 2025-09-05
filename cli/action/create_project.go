package action

import (
	generate2 "github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/tui/input"
)

func CreateProject(options CreateProjectOptions) (err error) {
	if options.Name == "" {
		options.Name, err = input.Send("give the project a name")
		if err != nil {
			return err
		}
	}

	return generate2.Project(generate2.ProjectOptions{
		Name: options.Name,
		Efs:  options.Efs,
	})
}
