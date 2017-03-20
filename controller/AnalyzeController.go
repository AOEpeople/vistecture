package controller

import (
	"fmt"
	"os"
	"appdependency/model/analyze"
)

type AnalyzeController struct {
	ProjectConfigPath string
}

func (AnalyzeController AnalyzeController) AnalyzeAction() {
	Project := loadProject(AnalyzeController.ProjectConfigPath)
	var ProjectAnalyzer analyze.ProjectAnalyzer
	e := ProjectAnalyzer.Analyze(Project)
	if e != nil {
		fmt.Println(e)
		os.Exit(-1)
	}
	fmt.Println("All right!")
}

func (AnalyzeController AnalyzeController) ValidateAction()  {
	loadProject(AnalyzeController.ProjectConfigPath)
	fmt.Println("Valid Project definition")
}
