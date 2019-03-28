package core

import (
	"errors"
	"html/template"
	"strings"

	"github.com/russross/blackfriday"
)

type (
	Application struct {
		//Id - generated Integer as Id
		Id                         int                        `json:"id" yaml:"id"`
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

	ApplicationDisplaySettings struct {
		Rotate      bool   `json:"rotate" yaml:"rotate"`
		BorderColor string `json:"color,omitempty" yaml:"borderColor,omitempty"`
		Color       string `json:"color,omitempty" yaml:"color,omitempty"`
	}

	//DependenciesGrouped Value object represents all Dependencies to one Application
	DependenciesGrouped struct {
		Application       *Application `json:"application"`
		SourceApplication *Application `json:"sourceApplication"`
		Dependencies      []Dependency `json:"dependencies"`
	}

	InfrastructureDependency struct {
		Type string `json:"type" yaml:"type"`
	}
)

const (
	STATUS_PLANNED    = "planned"
	CATEGORY_EXTERNAL = "external"
)

//Validate - validates the Application
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

//GetDescriptionHtml - helper that renders the description text as markdown - to be used in HTML documentations
func (a *Application) GetDescriptionHtml() template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(a.Description)))
}

//GetSummary -  returns summary. If summary is not set the first 100 letters from description
func (a *Application) GetSummary() string {
	if a.Summary != "" {
		return a.Summary
	}
	if len(a.Description) > 100 {
		return a.Description[0:100] + "..."
	}
	return a.Description
}

//FindService - returns service object or error
func (a *Application) FindService(nameToMatch string) (*Service, error) {
	for _, service := range a.ProvidedServices {
		if service.Name == nameToMatch {
			return &service, nil
		}
	}
	return nil, errors.New("a '" + a.Name + "' has no service with Name " + nameToMatch)
}

//returns the  Applications that are a dependency of the current application in the passed project
func (a *Application) GetAllDependencyApplications(Project *Project) ([]*Application, error) {
	var result []*Application
	for _, dependency := range a.GetAllDependencies() {
		foundComponent, e := dependency.GetApplication(Project)
		if e != nil {
			return nil, e
		}
		result = append(result, foundComponent)
	}
	return result, nil
}

//returns all the  Dependencies objects from the current application to others
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

//GetDependenciesTo returns the Dependencies to a specified other application
func (a *Application) GetDependenciesTo(applicationName string) ([]Dependency, error) {
	var result []Dependency
	for _, dependency := range a.GetAllDependencies() {
		if dependency.GetApplicationName() == applicationName {
			result = append(result, dependency)
		}
	}
	if len(result) == 0 {
		return nil, errors.New("Dependency to '" + applicationName + "' Not found")
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

//IsOpenHostApp - returns true if the APIs provided by this service are all declared as IsOpenHost.
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

//GetDependenciesGrouped - returns a list of grouped dependencies for this application to others. Useful if you are not interested in the indivudual dependencies but only the general "links" from this app to others
func (a *Application) GetDependenciesGrouped(project *Project) []*DependenciesGrouped {
	var result []*DependenciesGrouped

	//private helper func
	getGroupedDepToAppInList := func(list []*DependenciesGrouped, depApp *Application) *DependenciesGrouped {
		for _, group := range list {
			if group.Application == depApp {
				return group

			}
		}
		return nil
	}

	for _, dep := range a.Dependencies {
		depApp, err := dep.GetApplication(project)
		if err != nil {
			continue
		}
		//add dependency to existing depGrouped if found
		groupedDep := getGroupedDepToAppInList(result, depApp)
		if groupedDep != nil {
			groupedDep.Dependencies = append(groupedDep.Dependencies, dep)
			continue
		}
		//else append new
		result = append(result, &DependenciesGrouped{
			Application:       depApp,
			SourceApplication: a,
			Dependencies:      []Dependency{dep},
		})
	}
	return result
}

//GetMissingDependencies - returns a list of references application names that are not in the project
func (a *Application) GetMissingDependencies(project *Project) []string {
	var missing []string

	for _, dep := range a.Dependencies {
		_, err := dep.GetApplication(project)
		if err != nil {
			missing = append(missing, dep.GetApplicationName())
		}
	}
	return missing
}
