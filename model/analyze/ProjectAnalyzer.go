package analyze

import (
	"errors"
	"strconv"

	"github.com/AOEpeople/vistecture/model/core"
)

type ProjectAnalyzer struct{}

//Analyze validates project and Components
func (projectAnalyzer *ProjectAnalyzer) AnalyzeCyclicDependencies(project *core.Project) []error {
	//Walk dependencies add add service to called stack
	var stack []string
	var errors []error
	for _, component := range project.Applications {
		err := projectAnalyzer.walkDependencies(project, component, stack)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

//Analyze validates project and Components
func (projectAnalyzer *ProjectAnalyzer) ImpactAnalyze(project *core.Project) []string {
	//Walk dependencies add add service to called stack
	var impactsPerComponent []string
	for _, component := range project.Applications {
		directComponents := project.FindApplicationThatReferenceTo(component, false)
		allIndirectComponents := project.FindApplicationThatReferenceTo(component, true)
		impactsPerComponent = append(impactsPerComponent, strconv.Itoa(len(directComponents))+"\t\t"+strconv.Itoa(len(allIndirectComponents))+"\t\t"+component.Name)
	}
	return impactsPerComponent
}

// called recursive and adds the dependency components to the stack
// if a component appears again it throws an error (= cyclic dependencie)
func (projectAnalyzer *ProjectAnalyzer) walkDependencies(project *core.Project, component *core.Application, callStack []string) error {
	//end of recursion
	if len(callStack) >= 50 {
		return errors.New("More than 50 depth")
	}
	if arrayContainsName(callStack, component.Name) {
		errorMessage := "Cyclic dependency: "

		for _, v := range callStack {
			errorMessage += v + " -> "
		}
		return errors.New(errorMessage + component.Name)
	}

	callStack = append(callStack, component.Name)
	//Walk Dependcies that are global for the given component
	dependencies := component.Dependencies
	for _, dependency := range dependencies {
		nextComponent, _ := dependency.GetComponent(project)
		err := projectAnalyzer.walkDependencies(project, &nextComponent, callStack)
		if err != nil {
			return err
		}
	}
	//Walk Dependencies that are related to the called service
	service, e := component.FindService("ssss")
	if e == nil {
		dependencies2 := service.Dependencies
		for _, dependency := range dependencies2 {
			nextComponent, _ := dependency.GetComponent(project)
			err := projectAnalyzer.walkDependencies(project, &nextComponent, callStack)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func arrayContainsName(callStack []string, name string) bool {
	for _, v := range callStack {
		if v == name {
			return true
		}
	}
	return false
}
