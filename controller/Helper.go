package controller

import (
	"fmt"
	"os"
	"appdependency/model/core"
)


func loadProject(ProjectConfigPath string) *core.Project {
	project,e := core.CreateProjectAndValidate(ProjectConfigPath)
	if (e != nil) {
		fmt.Println("Project JSON is not valid:", e)
		os.Exit(-1)
	}
	return project
}