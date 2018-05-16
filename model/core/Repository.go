package core

import (
	"errors"
)

type Repository struct {
	ProjectConfig []*ProjectConfig `json:"projects" yaml:"projects"`
	Applications  []*Application   `json:"applications" yaml:"applications"`
}

//Validates repository
func (Repository *Repository) Validate() []error {
	var foundErrors []error

	for _, projectInfo := range Repository.ProjectConfig {
		foundErrors = append(foundErrors, projectInfo.Validate()...)
	}

	return foundErrors
}

//Find application by Name
func (Repository *Repository) FindApplicationByName(nameToMatch string) (Application, error) {
	for _, component := range Repository.Applications {
		if component.Name == nameToMatch {
			return *component, nil
		}
	}
	return Application{}, errors.New("Application with name '" + nameToMatch + "' not found")
}

func (Repository *Repository) FindApplicationsByCategory(categoryToMatch Category) []*Application {

	var resultApplications []*Application
	for _, currentApplications := range Repository.Applications {
		if currentApplications.Category == categoryToMatch.Value() {
			resultApplications = append(resultApplications, currentApplications)
		}
	}
	return resultApplications
}

//Check if a component with Name exist
func (Repository *Repository) HasApplicationWithName(nameToMatch string) bool {
	if _, e := Repository.FindApplicationByName(nameToMatch); e != nil {
		return false
	}
	return true
}

//Find project info by Name
func (Repository *Repository) FindProjectConfigByName(nameToMatch string) (ProjectConfig, error) {
	for _, projectConfig := range Repository.ProjectConfig {
		if projectConfig.Name == nameToMatch {
			return *projectConfig, nil
		}
	}
	return ProjectConfig{}, errors.New("project info with name '" + nameToMatch + "' not found")
}

//Gets the project info by name. If the name is not found, return the first available one.
func (Repository *Repository) GetProjectConfig(nameToMatch string) *ProjectConfig {
	projectConfig, error := Repository.FindProjectConfigByName(nameToMatch)
	if error != nil {
		if len(Repository.ProjectConfig) >= 1 {
			return Repository.ProjectConfig[0]

		} else {
			return &ProjectConfig{Name: "Full Repository"}
		}
	}
	return &projectConfig
}

//Check if a component with Name exist
func (Repository *Repository) HasProjectInfoWithName(nameToMatch string) bool {
	if _, e := Repository.FindProjectConfigByName(nameToMatch); e != nil {
		return false
	}
	return true
}

//Merges the given repository with another. The current repository is the one who will be modified.
func (Repository *Repository) MergeWith(OtherRepository *Repository) error {
	for _, application := range OtherRepository.Applications {
		if Repository.HasApplicationWithName(application.Name) {
			return errors.New("Application name: '" + application.Name + "' Is duplicated")
		}
		Repository.Applications = append(Repository.Applications, application)
	}

	for _, projectInfo := range OtherRepository.ProjectConfig {
		if Repository.HasProjectInfoWithName(projectInfo.Name) {
			return errors.New("project name: '" + projectInfo.Name + "' Is duplicated")
		}
		Repository.ProjectConfig = append(Repository.ProjectConfig, projectInfo)
	}
	return nil
}
