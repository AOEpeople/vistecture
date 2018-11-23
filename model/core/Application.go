package core

import (
	"errors"
	"html/template"
	"strings"

	"github.com/russross/blackfriday"
)

type (
	Application struct {
		Name                       string                     `json:"name" yaml:"name"`
		Team                       string                     `json:"team" yaml:"team"`
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
		Status                     string                     `json:"status" yaml:"status"`
	}

	//DependenciesGrouped Value object represents all Dependencies to one Application
	DependenciesGrouped struct {
		Application       *Application
		SourceApplication *Application
		Dependencies      []Dependency
	}
)

const (
	STATUS_PLANNED = "planned"
)

func (a *Application) Validate() []error {
	var foundErrors []error

	if len(a.Name) <= 0 {
		foundErrors = append(foundErrors, errors.New("a with no name found."))
	}
	if strings.Contains(a.Name, ".") {
		foundErrors = append(foundErrors, errors.New("a name contains '.'"))
	}
	return foundErrors
}

func (a *Application) GetDescriptionHtml() template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(a.Description)))
}

// Returns summary. If summary is not set the first 100 letters from description
func (a *Application) GetSummary() string {
	if a.Summary != "" {
		return a.Summary
	}
	if len(a.Description) > 100 {
		return a.Description[0:100] + "..."
	}
	return a.Description
}

func (a *Application) FindService(nameToMatch string) (*Service, error) {
	for _, service := range a.ProvidedServices {
		if service.Name == nameToMatch {
			return &service, nil
		}
	}

	return nil, errors.New("a '" + a.Name + "' has no Interface with Name " + nameToMatch)
}

//returns the  Applications that are a dependency of the current application
func (a *Application) GetAllDependencyApplications(Project *Project) ([]*Application, error) {
	var result []*Application
	// Walk dependencies from current component
	for _, dependency := range a.Dependencies {
		foundComponent, e := dependency.GetApplication(Project)
		if e != nil {
			return nil, e
		}
		result = append(result, foundComponent)
	}
	// Walk dependencies - modeled from current components provides Services
	for _, service := range a.ProvidedServices {
		for _, dependency := range service.Dependencies {
			foundComponent, e := dependency.GetApplication(Project)
			if e != nil {
				return nil, e
			}
			result = append(result, foundComponent)
		}
	}

	return result, nil
}

//returns the depending Dependencies
func (a *Application) GetDependenciesTo(ComponentName string) ([]Dependency, error) {
	var result []Dependency
	for _, dependency := range a.GetAllDependencies() {
		if dependency.GetApplicationName() == ComponentName {
			result = append(result, dependency)
		}
	}
	if len(result) == 0 {
		return nil, errors.New("Dependency to '" + ComponentName + "' Not found")
	}
	return result, nil
}

//GetServiceForDependency  - returns the provided service that is supposed to be referenced by the Dependency
func (a *Application) GetServiceForDependency(dependency *Dependency) *Service {
	if dependency.GetApplicationName() != a.Name {
		return nil
	}

	if dependency.GetServiceName() == "" {
		return nil
	}
	service, err := a.FindService(dependency.GetServiceName())
	if err != nil {
		return nil
	}
	return service
}

//returns the depending Dependencies
func (a *Application) GetAllDependencies() []Dependency {
	var result []Dependency
	for _, dependency := range a.Dependencies {
		result = append(result, dependency)
	}
	for _, service := range a.ProvidedServices {
		for _, dependency := range service.Dependencies {
			result = append(result, dependency)
		}
	}
	return result
}

func (a *Application) IsOpenHostApp() bool {
	if len(a.ProvidedServices) == 0 {
		return false
	}
	for _, service := range a.ProvidedServices {
		if !service.IsOpenHost && service.Type != "gui" {
			return false
		}
	}
	return true
}

//GetGroupPath - returns the list of Groups the application is part of (parent to leaf)
func (a *Application) GetGroupPath() []string {
	return strings.Split(a.Group, "/")
}

//GetDependenciesGrouped - returns a list of grouped dependencies by application
func (a *Application) GetDependenciesGrouped(project *Project) []*DependenciesGrouped {
	var result []*DependenciesGrouped
	for _, dep := range a.Dependencies {
		depApp, err := dep.GetApplication(project)
		if err != nil {
			continue
		}
		inList := false
		for _, group := range result {
			if group.Application == depApp {
				inList = true
				group.Dependencies = append(group.Dependencies, dep)
			}
		}
		if !inList {
			result = append(result, &DependenciesGrouped{
				Application:       depApp,
				SourceApplication: a,
				Dependencies:      []Dependency{dep},
			})
		}
	}
	return result
}

//Merges the given application with another. The current application is the one who will be modified.
func (a *Application) GetMerged(applicationReference ApplicationReference) (*Application, error) {
	newApplication := *a

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
		newApplication.Dependencies = append(a.Dependencies, applicationReference.AddDependencies...)
	}
	if applicationReference.AddProvidedServices != nil {
		newApplication.ProvidedServices = append(a.ProvidedServices, applicationReference.AddProvidedServices...)
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
	return &newApplication, nil
}
