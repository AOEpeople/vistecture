package core

import (
	"errors"
)

type Project struct {
	Name         string
	Applications []*Application
	// TODO - add services
}

const NOGROUP = "nogroup"

//Validates Project and Components
func (Project *Project) Validate() []error {

	var foundErrors []error

	for _, application := range Project.Applications {
		foundErrors = append(foundErrors, application.Validate()...)
		dependencies := application.GetAllDependencies()

		for _, dependency := range dependencies {
			dependendComponentName, serviceName := dependency.GetComponentAndServiceNames()

			error := Project.doesServiceExists(dependendComponentName, serviceName)
			if error != nil {
				foundErrors = append(foundErrors, errors.New("Application '" + application.Name + "' Dependencies has Error: " + error.Error()))
			}
		}
	}
	return foundErrors
}

//Find by Name
func (Project *Project) FindApplication(nameToMatch string) (Application, error) {
	for _, component := range Project.Applications {
		if component.Name == nameToMatch {
			return *component, nil
		}
	}
	return Application{}, errors.New("Application with name '" + nameToMatch + "' not found")
}

//Find by Name
func (Project *Project) FindApplicationThatReferenceTo(application *Application, recursive bool) []*Application {
	if recursive {
		// api.2 -> ma -> api.1
		return Project.findAllApplicationsThatReferenceComponent(application)
	} else {
		return Project.findApplicationsThatReferenceComponent(application)
	}

}

// Get Map with components grouped by Group.
// NOGROUP is used for ungrouped components
func (Project *Project) GetApplicationByGroup() map[string][]*Application {
	m := make(map[string][]*Application)
	for _, component := range Project.Applications {
		if len(component.Group) > 0 {
			m[component.Group] = append(m[component.Group], component)
		} else {
			m[NOGROUP] = append(m[NOGROUP], component)
		}
	}
	return m
}

func (Project *Project) findAllApplicationsThatReferenceComponent(componentReferenced *Application) []*Application {
	var referencingComponents []*Application
	for _, currentComponent := range Project.findApplicationsThatReferenceComponent(componentReferenced) {
		referencingComponents = append(referencingComponents, currentComponent)
		recursiveReferencingComponents := Project.findAllApplicationsThatReferenceComponent(currentComponent)
		for _, currentComponent := range recursiveReferencingComponents {
			referencingComponents = append(referencingComponents, currentComponent)
			if sliceContains(referencingComponents, currentComponent) {
				continue
			}
		}
	}
	return referencingComponents
}

// returns all components that have a direct dependency to the given component
func (Project *Project) findApplicationsThatReferenceComponent(componentReferenced *Application) []*Application {
	var referencingComponents []*Application
	//walk through all registered compoents and return those who match
	for _, currentComponent := range Project.Applications {
		currentComponentsDependencies := currentComponent.GetAllDependencies()
		for _, dependency := range currentComponentsDependencies {
			if dependency.GetComponentName() != componentReferenced.Name {
				continue
			}
			if sliceContains(referencingComponents, currentComponent) {
				continue
			}
			referencingComponents = append(referencingComponents, currentComponent)
		}
	}
	return referencingComponents
}

//Helper to check if a slice of components already contains a certain Component
func sliceContains(searchInList []*Application, searchFor *Application) bool {
	for _, s := range searchInList {
		if s == searchFor {
			return true
		}
	}
	return false
}

// internal method - Checks if a service exists and returns error if not
func (Project *Project) doesServiceExists(dependendComponentName string, serviceName string) error {
	dependendComponent, errorOnComponentFound := Project.FindApplication(dependendComponentName)
	if errorOnComponentFound != nil {
		return errorOnComponentFound
	}
	if serviceName == "" {
		return nil
	}
	if _, errorOnServiceFound := dependendComponent.FindService(serviceName); errorOnServiceFound != nil {
		return errorOnServiceFound
	}
	return nil
}
