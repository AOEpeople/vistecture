package controller

import (
	"log"

	"github.com/AOEpeople/vistecture/model/core"
)

func loadProject(ProjectConfigPath string, ProjectName string) *core.Project {
	project, errors := core.CreateProjectByName(ProjectConfigPath, ProjectName)

	if len(errors) > 0 {
		for _, err := range errors {
			log.Print("project creation failed because of: ", err)
		}
		log.Fatal("project loading aborted.")
	}
	return project
}
