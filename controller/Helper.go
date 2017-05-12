package controller

import (
	"log"

	"github.com/AOEpeople/vistecture/model/core"
)

func loadProject(ProjectConfigPath string) *core.Project {
	project, err := core.CreateProject(ProjectConfigPath)
	if err != nil {
		log.Fatal("Project JSON is not valid:", err)
	}
	err = project.Validate()
	if err != nil {
		log.Fatal("Validation Errors:", err)
	}
	return project
}
