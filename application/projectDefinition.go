package application

import (
	"errors"
	"strings"

	"github.com/AOEpeople/vistecture/model/core"
)

type (
	ProjectDefinitions struct {
		ProjectConfig []*ProjectConfig    `json:"projects" yaml:"projects"`
		Applications  []*core.Application `json:"applications" yaml:"applications"`
	}
	ProjectConfig struct {
		Name                string                  `json:"name" yaml:"name" `
		IncludedApplication []*ApplicationReference `json:"included-applications" yaml:"included-applications"`
	}
	ApplicationReference struct {
		//Name - is used to reference
		Name string `json:"name" yaml:"name"`
		//Title and all other attributes are supposed to override or extend the properties of the referenced application
		Title               string            `json:"title,omitempty" yaml:"title,omitempty"`
		Summary             string            `json:"summary,omitempty" yaml:"summary,omitempty"`
		Description         string            `json:"description,omitempty" yaml:"description,omitempty"`
		Group               string            `json:"group,omitempty" yaml:"group,omitempty"`
		Technology          string            `json:"technology,omitempty" yaml:"technology,omitempty"`
		Category            string            `json:"category,omitempty" yaml:"category,omitempty"`
		AddProvidedServices []core.Service    `json:"add-provided-services" yaml:"add-provided-services"`
		AddDependencies     []core.Dependency `json:"add-dependencies" yaml:"add-dependencies"`
		Properties          map[string]string `json:"properties" yaml:"properties"`
	}
)

//Validates project and Components
func (p *ProjectConfig) Validate() []error {
	var foundErrors []error

	if len(p.Name) <= 0 {
		foundErrors = append(foundErrors, errors.New("project config with no name found."))
	}
	if strings.Contains(p.Name, ".") {
		foundErrors = append(foundErrors, errors.New("project config name contains '.'"))
	}
	return foundErrors
}

//Validates Definitions
func (p *ProjectDefinitions) Validate() []error {
	var foundErrors []error

	for _, projectInfo := range p.ProjectConfig {
		foundErrors = append(foundErrors, projectInfo.Validate()...)
	}

	return foundErrors
}

//Find application by Name
func (p *ProjectDefinitions) FindApplicationByName(nameToMatch string) (*core.Application, error) {
	for _, component := range p.Applications {
		if component.Name == nameToMatch {
			return component, nil
		}
	}
	return nil, errors.New("Application with name '" + nameToMatch + "' not found")
}

//Gets the project info by name. If the name is not found, return the first available one.
func (p *ProjectDefinitions) GetProjectConfig(nameToMatch string) *ProjectConfig {
	projectConfig, error := p.findProjectConfigByName(nameToMatch)
	if error != nil {
		if len(p.ProjectConfig) >= 1 {
			return p.ProjectConfig[0]

		} else {
			return &ProjectConfig{Name: "Full Project Definitions"}
		}
	}
	return &projectConfig
}

//Check if a component with Name exist
func (p *ProjectDefinitions) hasApplicationWithName(nameToMatch string) bool {
	if _, e := p.FindApplicationByName(nameToMatch); e != nil {
		return false
	}
	return true
}

//Find project info by Name
func (Repository *ProjectDefinitions) findProjectConfigByName(nameToMatch string) (ProjectConfig, error) {
	for _, projectConfig := range Repository.ProjectConfig {
		if projectConfig.Name == nameToMatch {
			return *projectConfig, nil
		}
	}
	return ProjectConfig{}, errors.New("project info with name '" + nameToMatch + "' not found")
}

//Check if a component with Name exist
func (p *ProjectDefinitions) hasProjectInfoWithName(nameToMatch string) bool {
	if _, e := p.findProjectConfigByName(nameToMatch); e != nil {
		return false
	}
	return true
}

//Merges the given repository with another. The current repository is the one who will be modified.
func (p *ProjectDefinitions) mergeWith(otherRepository *ProjectDefinitions) error {
	if otherRepository == nil {
		return errors.New("No OtherRepository given")
	}
	for _, application := range otherRepository.Applications {
		if p.hasApplicationWithName(application.Name) {
			return errors.New("Application name: '" + application.Name + "' Is duplicated")
		}
		p.Applications = append(p.Applications, application)
	}

	for _, projectInfo := range otherRepository.ProjectConfig {
		if p.hasProjectInfoWithName(projectInfo.Name) {
			return errors.New("project name: '" + projectInfo.Name + "' Is duplicated")
		}
		p.ProjectConfig = append(p.ProjectConfig, projectInfo)
	}
	return nil
}

//Merges the given application with another. The current application is the one who will be modified.
func (a *ApplicationReference) GetAdjustedApplication(application *core.Application) (*core.Application, error) {
	//clone the passed application
	newApplication := *application

	//override properties if set in the ApplicationReference
	if a.Name != "" {
		newApplication.Name = a.Name
	}
	if a.Category != "" {
		newApplication.Category = a.Category
	}
	if a.Description != "" {
		newApplication.Description = a.Description
	}
	if a.AddDependencies != nil {
		newApplication.Dependencies = append(application.Dependencies, a.AddDependencies...)
	}
	if a.AddProvidedServices != nil {
		newApplication.ProvidedServices = append(application.ProvidedServices, a.AddProvidedServices...)
	}
	if a.Category != "" {
		newApplication.Category = a.Category
	}
	if a.Group != "" {
		newApplication.Group = a.Group
	}
	if a.Properties != nil {
		if newApplication.Properties == nil {
			newApplication.Properties = make(map[string]string)
		}
		for k, v := range a.Properties {
			newApplication.Properties[k] = v
		}
	}
	return &newApplication, nil
}
