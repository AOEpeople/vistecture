package controller

import (
	"log"

	"github.com/danielpoe/appdependency/model/core"
)

func loadProject(ProjectConfigPath string) *core.Project {
	project, err := core.CreateProjectAndValidate(ProjectConfigPath)
	if err != nil {
		log.Fatal("Project JSON is not valid:", err)
	}
	return project
}
