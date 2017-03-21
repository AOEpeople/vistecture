package analyze

import (
	"errors"

	"appdependency/model/core"
)

type ProjectAnalyzer struct{}

//Analyze validates Project and Components
func (projectAnalyzer *ProjectAnalyzer) Analyze(project *core.Project) error {
	//Walk dependencies add add service to called stack
	var stack []string
	for _, component := range project.Components {
		err := projectAnalyzer.walkDependencies(project, component, stack)
		if err != nil {
			return err
		}
	}
	return nil
}

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

	for _, dependency := range component.Dependencies {
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
