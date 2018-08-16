package controller

import (
	"log"

	"github.com/AOEpeople/vistecture/model/core"
)

func loadProject(ProjectConfigPath string, ProjectName string, skipValidation bool) *core.Project {
	project, errors := core.CreateProjectByName(ProjectConfigPath, ProjectName, skipValidation)

	if len(errors) > 0 {
		for _, err := range errors {
			log.Print("project creation failed because of: ", err)
		}
		log.Fatal("project loading aborted.")
	}
	return project
}
