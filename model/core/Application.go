package core

import (
	"errors"
	"fmt"
	"strings"

	"html/template"

	"github.com/russross/blackfriday"
)

type Application struct {
	Name                       string                     `json:"name" yaml:"name"`
	Title                      string                     `json:"title" yaml:"title"`
	Summary                    string                     `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description                string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Group                      string                     `json:"group,omitempty" yaml:"group,omitempty"`
	Technology                 string                     `json:"technology,omitempty" yaml:"technology,omitempty"`
	Category                   string                     `json:"category,omitempty" yaml:"category,omitempty"`
	ProvidedServices           []Service                  `json:"provided-services" yaml:"provided-services"`
	InfrastructureDependencies []InfrastructureDependency `json:"infrastructure-dependencies" yaml:"infrastructure-dependencies"`
	Dependencies               []Dependency               `json:"dependencies" yaml:"dependencies"`
	Display                    ApplicationDisplaySettings `json:"display,omitempty" yaml:"display,omitempty"`
	Properties                 map[string]string          `json:"properties,omitempty" yaml:"properties,omitempty"`
}

func (Component Application) Validate() bool {
	if strings.Contains(Component.Name, ".") {
		fmt.Printf("Component Name contains . '%v'\n", Component.Name)
		return false
	}
	return true
}

func (Component Application) GetDescriptionHtml() template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(Component.Description)))
}

// Returns summary. If summary is not set the first 100 letters from description
func (Component Application) GetSummary() string {
	if Component.Summary != "" {
		return Component.Summary
	}
	if len(Component.Description) > 100 {
		return Component.Description[0:100] + "..."
	}
	return Component.Description
}

func (Component Application) FindService(nameToMatch string) (*Service, error) {
	for _, service := range Component.ProvidedServices {
		if service.Name == nameToMatch {
			return &service, nil
		}
	}

	return nil, errors.New("Component " + Component.Name + " has no Interface with Name " + nameToMatch)
}

//returns the  Applications that are a dependency of the current application
func (GivenComponent Application) GetAllDependencyApplications(Project *Project) ([]Application, error) {
	var result []Application
	// Walk dependencies from current component
	for _, dependency := range GivenComponent.Dependencies {
		foundComponent, e := dependency.GetComponent(Project)
		if e != nil {
			return nil, e
		}
		result = append(result, foundComponent)
	}
	// Walk dependencies - modeled from current components provides Services
	for _, service := range GivenComponent.ProvidedServices {
		for _, dependency := range service.Dependencies {
			foundComponent, e := dependency.GetComponent(Project)
			if e != nil {
				return nil, e
			}
			result = append(result, foundComponent)
		}
	}

	return result, nil
}

//returns the depending Dependencies
func (GivenComponent Application) GetDependencyTo(ComponentName string) (Dependency, error) {
	var emptyDependency Dependency
	for _, dependency := range GivenComponent.GetAllDependencies() {
		if dependency.GetComponentName() == ComponentName {
			return dependency, nil
		}
	}
	return emptyDependency, errors.New("Dependency to " + ComponentName + " Not found")
}

//returns the depending Dependencies
func (GivenComponent Application) GetAllDependencies() []Dependency {
	var result []Dependency
	for _, dependency := range GivenComponent.Dependencies {
		result = append(result, dependency)
	}
	for _, service := range GivenComponent.ProvidedServices {
		for _, dependency := range service.Dependencies {
			result = append(result, dependency)
		}
	}
	return result
}
