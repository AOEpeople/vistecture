package core

import (
	"errors"
)

type Project struct {
	Name       string       `json:"name"`
	Components []*Component `json:"components"`
}

const NOGROUP = "nogroup"

//Validates Project and Components
func (Project *Project) Validate() error {
	for _, component := range Project.Components {
		if component.Validate() == false {
			return errors.New("Component not valid")
		}
		for _, dependency := range component.Dependencies {
			dependendComponentName, serviceName := dependency.GetComponentAndServiceNames()
			dependendComponent, errorOnComponentFound := Project.FindComponent(dependendComponentName)
			if errorOnComponentFound != nil {
				return errors.New("Component " + component.Name + " Dependencies has Error:" + errorOnComponentFound.Error())
			}
			if serviceName == "" {
				continue
			}
			if _, errorOnServiceFound := dependendComponent.FindService(serviceName); errorOnServiceFound != nil {
				return errorOnServiceFound
			}
		}
	}
	return nil
}

//Find by Name
func (Project *Project) FindComponent(nameToMatch string) (Component, error) {
	for _, component := range Project.Components {
		if component.Name == nameToMatch {
			return *component, nil
		}
	}
	return Component{}, errors.New("Component with name " + nameToMatch + " not found")
}

//Check if a component with Name exist
func (Project *Project) HasComponentWithName(nameToMatch string) bool {
	if _, e := Project.FindComponent(nameToMatch); e != nil {
		return false
	}
	return true
}

//Get Map with components grouped by Group.
// NOGROUP is used for ungrouped components
func (Project *Project) GetComponentsByGroup() map[string][]*Component {
	m := make(map[string][]*Component)
	for _, component := range Project.Components {
		if len(component.Group) > 0 {
			m[component.Group] = append(m[component.Group], component)
		} else {
			m[NOGROUP] = append(m[NOGROUP], component)
		}
	}
	return m
}

//Merges the given Project with another
func (Project *Project) AddComponentsFromProject(OtherProject *Project) error {
	for _, component := range OtherProject.Components {
		if Project.HasComponentWithName(component.Name) {
			return errors.New(component.Name + " Is duplicated")
		}
		Project.Components = append(Project.Components, component)
	}
	return nil
}
