package controller

import (
	"fmt"
	"log"
	"appdependency/model/analyze"
)

type AnalyzeController struct {
	ProjectConfigPath *string
}

func (AnalyzeController AnalyzeController) AnalyzeAction() {
	var project = loadProject(*AnalyzeController.ProjectConfigPath)
	var ProjectAnalyzer analyze.ProjectAnalyzer
	err := ProjectAnalyzer.Analyze(project)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All right!")
}

func (AnalyzeController AnalyzeController) ValidateAction() {
	loadProject(*AnalyzeController.ProjectConfigPath)
	fmt.Println("Valid Project definition")
}
