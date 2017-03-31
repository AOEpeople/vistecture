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
	errors := ProjectAnalyzer.AnalyzeCyclicDependencies(project)
	if errors != nil {
		for _,error := range errors {
			fmt.Println(error)
		}
		log.Fatal("Solve Errors please!")
	}

	fmt.Println()
	fmt.Println("Impact Analysis")
	impacts := ProjectAnalyzer.ImpactAnalyze(project)
	fmt.Println("Direct\t\tIndirect\tComponent")
	fmt.Println("------\t\t--------\t--------")
	for _,impact := range impacts {
		fmt.Println(impact)
	}

	fmt.Println("No errors found in your definitions!")


}

func (AnalyzeController AnalyzeController) ValidateAction() {
	loadProject(*AnalyzeController.ProjectConfigPath)
	fmt.Println("Valid Project definition")
}
