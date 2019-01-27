package core

import (
	"errors"
)

type (
	Project struct {
		Name         string         `json:"name" yaml:"name"`
		Applications []*Application `json:"applications" yaml:"applications"`
	}

	ApplicationsByGroup struct {
		SubGroups    []*ApplicationsByGroup `json:"subGroups"`
		Applications []*Application         `json:"applications"`
		GroupName    string                 `json:"groupName"`
		IsRoot       bool                   `json:"isRoot"`
	}
)

const (
	NOGROUP = "nogroup"
	NOTEAM  = "noteam"
)

//Validates project and Components
func (p *Project) Validate() []error {

	var foundErrors []error

	for _, application := range p.Applications {
		foundErrors = append(foundErrors, application.Validate()...)
		dependencies := application.GetAllDependencies()

		for _, dependency := range dependencies {
			dependendComponentName, serviceName := dependency.GetApplicationAndServiceNames()

			error := p.doesServiceExists(dependendComponentName, serviceName)
			if error != nil {
				foundErrors = append(foundErrors, errors.New("Application '"+application.Name+"' Dependencies has Error: "+error.Error()))
			}
		}
	}
	return foundErrors
}

func (p *Project) GenerateApplicationIds() {
	i := 1
	for _, app := range p.Applications {
		if app.Id != 0 {
			continue
		}
		app.Id = i
		i++
	}
}

//FindApplication - Find by Name
func (p *Project) FindApplication(nameToMatch string) (*Application, error) {
	for _, component := range p.Applications {
		if component.Name == nameToMatch {
			return component, nil
		}
	}
	return nil, errors.New("Application with name '" + nameToMatch + "' not found")
}

// GetApplicationsRootGroup - Returns the Root Group
func (p *Project) GetApplicationsRootGroup() *ApplicationsByGroup {
	appsByGroup := ApplicationsByGroup{
		IsRoot: true,
	}
	for _, app := range p.Applications {
		appsByGroup.add(app)
	}
	return &appsByGroup
}

// GetApplicationByTeam - Get Map with components grouped by Group.
// NOGROUP is used for ungrouped components
func (p *Project) GetApplicationByTeam() map[string][]*Application {
	m := make(map[string][]*Application)
	for _, component := range p.Applications {
		if len(component.Team) > 0 {
			m[component.Team] = append(m[component.Team], component)
		} else {
			m[NOTEAM] = append(m[NOTEAM], component)
		}
	}
	return m
}

//Find by Name
func (p *Project) FindApplicationThatReferenceTo(application *Application, recursive bool) []*Application {
	if recursive {
		// api.2 -> ma -> api.1
		return p.FindAllApplicationsThatReferenceApplication(application)
	} else {
		return p.FindApplicationsThatReferenceApplication(application)
	}

}

func (p *Project) FindAllApplicationsThatReferenceApplication(referencedApplication *Application) []*Application {
	var referencingApps []*Application
	for _, currentComponent := range p.FindApplicationsThatReferenceApplication(referencedApplication) {
		referencingApps = append(referencingApps, currentComponent)
		recursiveReferencingComponents := p.FindAllApplicationsThatReferenceApplication(currentComponent)
		for _, currentComponent := range recursiveReferencingComponents {
			referencingApps = append(referencingApps, currentComponent)
			if sliceContains(referencingApps, currentComponent) {
				continue
			}
		}
	}
	return referencingApps
}

// returns all components that have a direct dependency to the given component
func (p *Project) FindApplicationsThatReferenceApplication(referencedApplication *Application) []*Application {
	var referencingApps []*Application
	//walk through all registered compoents and return those who match
	for _, currentComponent := range p.Applications {
		currentComponentsDependencies := currentComponent.GetAllDependencies()
		for _, dependency := range currentComponentsDependencies {
			if dependency.GetApplicationName() != referencedApplication.Name {
				continue
			}
			if sliceContains(referencingApps, currentComponent) {
				continue
			}
			referencingApps = append(referencingApps, currentComponent)
		}
	}
	return referencingApps
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
func (p *Project) doesServiceExists(dependendComponentName string, serviceName string) error {
	dependendComponent, errorOnComponentFound := p.FindApplication(dependendComponentName)
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
