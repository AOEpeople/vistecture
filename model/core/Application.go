package core

import (
	"errors"
	"html/template"
	"strings"

	"github.com/russross/blackfriday"
)

type Application struct {
	Name                       string                     `json:"name" yaml:"name"`
	Title                      string                     `json:"title,omitempty" yaml:"title,omitempty"`
	Summary                    string                     `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description                string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Group                      string                     `json:"group,omitempty" yaml:"group,omitempty"`
	Technology                 string                     `json:"technology,omitempty" yaml:"technology,omitempty"`
	Category                   string                     `json:"category,omitempty" yaml:"category,omitempty"`
	ProvidedServices           []Service                  `json:"provided-services" yaml:"provided-services"`
	InfrastructureDependencies []InfrastructureDependency `json:"infrastructure-dependencies" yaml:"infrastructure-dependencies"`
	Dependencies               []Dependency               `json:"dependencies" yaml:"dependencies"`
	Display                    ApplicationDisplaySettings `json:"display,omitempty" yaml:"display,omitempty"`
	Properties                 map[string]string          `json:"properties" yaml:"properties"`
}

func (Application Application) Validate() []error {
	var foundErrors []error

	if len(Application.Name) <= 0 {
		foundErrors = append(foundErrors, errors.New("Application with no name found."))
	}
	if strings.Contains(Application.Name, ".") {
		foundErrors = append(foundErrors, errors.New("Application name contains '.'"))
	}
	return foundErrors
}

func (Application Application) GetDescriptionHtml() template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(Application.Description)))
}

// Returns summary. If summary is not set the first 100 letters from description
func (Application Application) GetSummary() string {
	if Application.Summary != "" {
		return Application.Summary
	}
	if len(Application.Description) > 100 {
		return Application.Description[0:100] + "..."
	}
	return Application.Description
}

func (Application Application) FindService(nameToMatch string) (*Service, error) {
	for _, service := range Application.ProvidedServices {
		if service.Name == nameToMatch {
			return &service, nil
		}
	}

	return nil, errors.New("Application '" + Application.Name + "' has no Interface with Name " + nameToMatch)
}

//returns the  Applications that are a dependency of the current application
func (GivenApplication Application) GetAllDependencyApplications(Project *Project) ([]Application, error) {
	var result []Application
	// Walk dependencies from current component
	for _, dependency := range GivenApplication.Dependencies {
		foundComponent, e := dependency.GetComponent(Project)
		if e != nil {
			return nil, e
		}
		result = append(result, foundComponent)
	}
	// Walk dependencies - modeled from current components provides Services
	for _, service := range GivenApplication.ProvidedServices {
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
func (GivenApplication Application) GetDependencyTo(ComponentName string) (Dependency, error) {
	var emptyDependency Dependency
	for _, dependency := range GivenApplication.GetAllDependencies() {
		if dependency.GetComponentName() == ComponentName {
			return dependency, nil
		}
	}
	return emptyDependency, errors.New("Dependency to '" + ComponentName + "' Not found")
}

//returns the depending Dependencies
func (GivenApplication *Application) GetAllDependencies() []Dependency {
	var result []Dependency
	for _, dependency := range GivenApplication.Dependencies {
		result = append(result, dependency)
	}
	for _, service := range GivenApplication.ProvidedServices {
		for _, dependency := range service.Dependencies {
			result = append(result, dependency)
		}
	}
	return result
}

//Merges the given application with another. The current application is the one who will be modified.
func (app *Application) GetMerged(applicationReference ApplicationReference) (Application, error) {
	newApplication := *app

	if applicationReference.Name != "" {
		newApplication.Name = applicationReference.Name
	}
	if applicationReference.Category != "" {
		newApplication.Category = applicationReference.Category
	}
	if applicationReference.Description != "" {
		newApplication.Description = applicationReference.Description
	}
	if applicationReference.AddDependencies != nil {
		newApplication.Dependencies = append(app.Dependencies, applicationReference.AddDependencies...)
	}
	if applicationReference.AddProvidedServices != nil {
		newApplication.ProvidedServices = append(app.ProvidedServices, applicationReference.AddProvidedServices...)
	}
	if applicationReference.Category != "" {
		newApplication.Category = applicationReference.Category
	}
	if applicationReference.Group != "" {
		newApplication.Group = applicationReference.Group
	}
	if applicationReference.Properties != nil {
		if newApplication.Properties == nil {
			newApplication.Properties = make(map[string]string)
		}
		for k, v := range applicationReference.Properties {
			newApplication.Properties[k] = v
		}
	}
	return newApplication, nil
}
