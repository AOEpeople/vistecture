package core

import (
	"errors"
	"strings"
)

type ProjectConfig struct {
	Name                string                  `json:"name" yaml:"name" `
	IncludedApplication []*ApplicationReference `json:"included-applications" yaml:"included-applications"`
}

//Validates project and Components
func (ProjectConfig *ProjectConfig) Validate() []error {
	var foundErrors []error

	if len(ProjectConfig.Name) <= 0 {
		foundErrors = append(foundErrors, errors.New("project config with no name found."))
	}
	if strings.Contains(ProjectConfig.Name, ".") {
		foundErrors = append(foundErrors, errors.New("project config name contains '.'"))
	}
	return foundErrors
}
