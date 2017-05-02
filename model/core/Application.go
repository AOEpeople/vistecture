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
	Description                string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Summary                    string                     `json:"summary,omitempty" yaml:"summary,omitempty"`
	Group                      string                     `json:"group,omitempty" yaml:"group,omitempty"`
	Technology                 string                     `json:"technology,omitempty" yaml:"technology,omitempty"`
	Category                   string                     `json:"category,omitempty" yaml:"category,omitempty"`
	ProvidedServices           []Service                  `json:"provided-services" yaml:"provided-services"`
	InfrastructureDependencies []InfrastructureDependency `json:"infrastructure-dependencies" yaml:"infrastructure-dependencies"`
	Dependencies               []Dependency               `json:"dependencies" yaml:"dependencies"`
	Display                    ApplicationDisplaySettings `json:"display,omitempty" yaml:"display,omitempty"`
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

//returns the depending Components
func (GivenComponent Application) GetAllRelatedComponents(Project *Project) ([]Application, error) {
	var result []Application
	for _, dependency := range GivenComponent.Dependencies {
		foundComponent, e := dependency.GetComponent(Project)
		if e != nil {
			return nil, e
		}
		result = append(result, foundComponent)
	}
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
func (GivenComponent Application) GetAllDependencies(Project *Project) ([]Dependency, error) {
	var result []Dependency
	for _, dependency := range GivenComponent.Dependencies {
		result = append(result, dependency)
	}
	for _, service := range GivenComponent.ProvidedServices {
		for _, dependency := range service.Dependencies {
			result = append(result, dependency)
		}
	}
	return result, nil
}
