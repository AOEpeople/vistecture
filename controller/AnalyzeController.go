package controller

import (
	"fmt"
	"log"

	"github.com/AOEpeople/vistecture/v2/model/analyze"
	"github.com/AOEpeople/vistecture/v2/model/core"
)

type AnalyzeController struct {
	project *core.Project
}

func (a *AnalyzeController) Inject(project *core.Project) {
	a.project = project
}

func (a *AnalyzeController) AnalyzeAction() {
	var ProjectAnalyzer analyze.ProjectAnalyzer
	errors := ProjectAnalyzer.AnalyzeCyclicDependencies(a.project)
	if errors != nil {
		for _, error := range errors {
			log.Println(error)
		}
		log.Fatal("Solve Errors please!")
	}
	fmt.Println("\nGreat - no errors or cyclic dependencies found in your definitions!")

	fmt.Println()
	fmt.Println("Impact Analysis: \n(How many other components may be influenced if a component fails)")
	impacts := ProjectAnalyzer.ImpactAnalyze(a.project)
	fmt.Println("Direct\t\tIndirect\tComponent")
	fmt.Println("------\t\t--------\t--------")
	for _, impact := range impacts {
		fmt.Println(impact)
	}

}
