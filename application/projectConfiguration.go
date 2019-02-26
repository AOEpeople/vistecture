package application

import (
	"errors"
	"strings"

	"github.com/AOEpeople/vistecture/v2/model/core"
)

type (
	ProjectConfig struct {
		SubViewConfig       []*SubViewConfig        `json:"subViews" yaml:"subViews"`
		AppDefinitionsPaths []string                `json:"appDefinitionsPaths" yaml:"appDefinitionsPaths"`
		ProjectName         string                  `json:"projectName" yaml:"projectName"`
		AppOverrides        []*ApplicationOverrides `json:"appOverrides" yaml:"appOverrides"`
	}
	SubViewConfig struct {
		Name                string   `json:"name" yaml:"name" `
		IncludedApplication []string `json:"included-applications" yaml:"included-applications"`
	}
	ApplicationOverrides struct {
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
func (p *SubViewConfig) validate() []error {
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
func (p *ProjectConfig) Validate() []error {
	var foundErrors []error

	for _, subView := range p.SubViewConfig {
		foundErrors = append(foundErrors, subView.validate()...)
	}

	return foundErrors
}

//Find project info by Name
func (p *ProjectConfig) FindSubViewConfigByName(nameToMatch string) (*SubViewConfig, error) {
	for _, subView := range p.SubViewConfig {
		if subView.Name == nameToMatch {
			return subView, nil
		}
	}
	return nil, errors.New("project info with name '" + nameToMatch + "' not found")
}

//Merges the given application with another. The current application is the one who will be modified.
func (a *ApplicationOverrides) GetAdjustedApplication(application *core.Application) (*core.Application, error) {
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

func (s *SubViewConfig) GetMatchedApps(apps []*core.Application) []*core.Application {
	var matchingApps []*core.Application
	for _, app := range apps {
		for _, includedAppName := range s.IncludedApplication {
			if includedAppName == app.Name {
				matchingApps = append(matchingApps, app)
			}
		}
	}
	return matchingApps
}
