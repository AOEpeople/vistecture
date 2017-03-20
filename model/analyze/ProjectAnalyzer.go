package analyze

import (
	"errors"
	core "appdependency/model/core"
)

type ProjectAnalyzer struct {
}


//Validates Project and Components
func (ProjectAnalyzer *ProjectAnalyzer) Analyze(Project *core.Project) error {

	//Walk dependencies add add service to called stack
	var stack []string
	for _, component := range Project.Components {
		e := ProjectAnalyzer.walkDependencies(Project,component,stack)
		if (e != nil) {
			return e
		}
	}
	return nil
}

func (ProjectAnalyzer *ProjectAnalyzer) walkDependencies(Project *core.Project,Component *core.Component, callStack []string) error {
	//end of recursion
	if len(callStack) >= 50 {
		return errors.New("More than 50 depth")
	}
	if (arrayContainsName(callStack,Component.Name)) {
		errorMessage := "Cyclic dependency: "

		for _,v := range callStack {
			errorMessage = errorMessage + v + " -> ";
		}
		return errors.New(errorMessage+ Component.Name)
	}

	//For next recursion copy map
	newCallStack := append(callStack, Component.Name)


	for _, dependency := range Component.Dependencies {
		nextComponent, _ := dependency.GetComponent(Project)
		e := ProjectAnalyzer.walkDependencies(Project,&nextComponent,newCallStack)
		if (e != nil) {
			return e
		}
	}
	return nil
}

func arrayContainsName(callStack []string, name string)  bool {
	for _, v := range callStack {
		if v == name {
			return true
		}
	}
	return false
}
