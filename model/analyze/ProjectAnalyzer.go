package analyze

import (
	"errors"

	"appdependency/model/core"
	"strconv"
)

type ProjectAnalyzer struct{}

//Analyze validates Project and Components
func (projectAnalyzer *ProjectAnalyzer) AnalyzeCyclicDependencies(project *core.Project) []error {
	//Walk dependencies add add service to called stack
	var stack []string
	var errors []error
	for _, component := range project.Components {
		err := projectAnalyzer.walkDependencies(project, component, stack)
		if err != nil {
			errors = append(errors,err)
		}
	}
	return errors
}

//Analyze validates Project and Components
func (projectAnalyzer *ProjectAnalyzer) ImpactAnalyze(project *core.Project) []string {
	//Walk dependencies add add service to called stack
	var impactsPerComponent []string
	for _, component := range project.Components {
		directComponents := project.FindComponentsThatReferenceTo(component, false)
		allIndirectComponents := project.FindComponentsThatReferenceTo(component, true)
		impactsPerComponent = append(impactsPerComponent,strconv.Itoa(len(directComponents)) + "\t\t"+ strconv.Itoa(len(allIndirectComponents)) +"\t\t" + component.Name)
	}
	return impactsPerComponent
}




// called recursive and adds the dependency components to the stack
// if a component appears again it throws an error (= cyclic dependencie)
func (projectAnalyzer *ProjectAnalyzer) walkDependencies(project *core.Project, component *core.Component, callStack []string) error {
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
	dependencies,_ := component.GetAllDependencies(project)
	for _,dependency  := range dependencies {
		nextComponent, _ := dependency.GetComponent(project)
		err := projectAnalyzer.walkDependencies(project, &nextComponent, callStack)
		if err != nil {
			return err
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
