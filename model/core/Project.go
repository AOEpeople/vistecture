package core

import (
	"errors"
)

type (
	Project struct {
		Name         string
		Applications []*Application
		// TODO - add services
	}

	ApplicationsByGroup struct {
		SubGroups    []*ApplicationsByGroup
		Applications []*Application
		GroupName    string
		IsRoot       bool
	}
)

const NOGROUP = "nogroup"
const NOTEAM = "noteam"

//Validates project and Components
func (Project *Project) Validate() []error {

	var foundErrors []error

	for _, application := range Project.Applications {
		foundErrors = append(foundErrors, application.Validate()...)
		dependencies := application.GetAllDependencies()

		for _, dependency := range dependencies {
			dependendComponentName, serviceName := dependency.GetComponentAndServiceNames()

			error := Project.doesServiceExists(dependendComponentName, serviceName)
			if error != nil {
				foundErrors = append(foundErrors, errors.New("Application '"+application.Name+"' Dependencies has Error: "+error.Error()))
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

// GetApplicationsRootGroup - Returns the Root Group
func (Project *Project) GetApplicationsRootGroup() *ApplicationsByGroup {
	appsByGroup := ApplicationsByGroup{
		IsRoot: true,
	}
	for _, app := range Project.Applications {
		appsByGroup.add(app)
	}
	return &appsByGroup
}

// GetApplicationByTeam - Get Map with components grouped by Group.
// NOGROUP is used for ungrouped components
func (Project *Project) GetApplicationByTeam() map[string][]*Application {
	m := make(map[string][]*Application)
	for _, component := range Project.Applications {
		if len(component.Team) > 0 {
			m[component.Team] = append(m[component.Team], component)
		} else {
			m[NOTEAM] = append(m[NOTEAM], component)
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

// Add a Application to the Value object and takes care that it is added to the correct subgroup
func (a *ApplicationsByGroup) add(app *Application) error {
	if !a.IsRoot {
		return errors.New("Only add to route group")
	}
	groupPath := app.GetGroupPath()
	groupToAddApp := a
	for _, group := range groupPath {
		groupToAddApp = groupToAddApp.getSubgroup(group)
	}
	groupToAddApp.Applications = append(groupToAddApp.Applications, app)
	return nil
}

// Add a Application to the Value object and takes care that it is added to the correct subgroup
func (a *ApplicationsByGroup) getSubgroup(groupName string) *ApplicationsByGroup {
	for _, subGroup := range a.SubGroups {
		if subGroup.GroupName == groupName {
			return subGroup
		}
	}
	newGroup := ApplicationsByGroup{
		GroupName: groupName,
	}
	a.SubGroups = append(a.SubGroups, &newGroup)
	return &newGroup
}
