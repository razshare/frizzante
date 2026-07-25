package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/razshare/frizzante/v2/cli/generations"
	"github.com/razshare/frizzante/v2/tui/inputs"
	"github.com/razshare/frizzante/v2/tui/messages"
)

func CreateProject(options CreateProjectOptions) (err error) {
	projectName := options.Name
	if projectName == "" {
		if options.Strict {
			err = errors.New("project name is missing")
			return
		}
		if projectName, err = inputs.Send("give the project a name"); err != nil {
			return
		}
	}
	if err = generations.Project(generations.ProjectOptions{
		Name: projectName,
		Efs:  options.Efs,
	}); err != nil {
		return
	}
	messages.Successf("project %s created with success!", projectName)
	step1 := fmt.Sprintf("1. cd %s", projectName)
	step2 := fmt.Sprintf("2. frizzante configure")
	step3 := fmt.Sprintf("3. frizzante dev")
	length1 := len(step1)
	length2 := len(step2)
	length3 := len(step3)
	width := length2
	if length1 > length2 {
		width = length1
	}
	padding1 := strings.Repeat(" ", width-length1)
	padding2 := strings.Repeat(" ", width-length2)
	padding3 := strings.Repeat(" ", width-length3)
	comment1 := fmt.Sprintf("%s#changes directory to %s", padding1, projectName)
	comment2 := fmt.Sprintf("%s#configures the project", padding2)
	comment3 := fmt.Sprintf("%s#starts development mode", padding3)
	messages.Tip(strings.Join([]string{
		"## next steps",
		fmt.Sprintf("%s %s", step1, comment1),
		fmt.Sprintf("%s %s", step2, comment2),
		fmt.Sprintf("%s %s", step3, comment3),
	}, "\n"))
	return
}
