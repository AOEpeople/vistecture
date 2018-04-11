package controller

import (
	"fmt"
	"log"

	"github.com/AOEpeople/vistecture/model/analyze"
)

type AnalyzeController struct {
	ProjectConfigPath *string
	ProjectName *string
}

func (AnalyzeController AnalyzeController) AnalyzeAction() {
	var project = loadProject(*AnalyzeController.ProjectConfigPath, *AnalyzeController.ProjectName)
	var ProjectAnalyzer analyze.ProjectAnalyzer
	errors := ProjectAnalyzer.AnalyzeCyclicDependencies(project)
	if errors != nil {
		for _, error := range errors {
			fmt.Println(error)
		}
		log.Fatal("Solve Errors please!")
	}
	fmt.Println("\nGreat - no errors or cyclic dependencies found in your definitions!")

	fmt.Println()
	fmt.Println("Impact Analysis: \n(How many other components may be influenced if a component fails)\n")
	impacts := ProjectAnalyzer.ImpactAnalyze(project)
	fmt.Println("Direct\t\tIndirect\tComponent")
	fmt.Println("------\t\t--------\t--------")
	for _, impact := range impacts {
		fmt.Println(impact)
	}

}

func (AnalyzeController AnalyzeController) ValidateAction() {
	loadProject(*AnalyzeController.ProjectConfigPath, *AnalyzeController.ProjectName)
	fmt.Println("Valid Project definition")
}
