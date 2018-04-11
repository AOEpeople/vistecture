package controller

import (
	"log"

	"github.com/AOEpeople/vistecture/model/core"
)

func loadProject(ProjectConfigPath string, ProjectName string) *core.Project {
	project, errors := core.CreateProjectByName(ProjectConfigPath, ProjectName)

	if len(errors) > 0 {
		for _, err := range errors {
			log.Print("Project creation failed because of: ", err)
		}
		log.Fatal("Project loading aborted.")
	}
	return project
}
