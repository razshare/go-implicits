package action

import (
	"github.com/razshare/go-implicits/internal/cli/generate"
	"github.com/razshare/go-implicits/tui/input"
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
