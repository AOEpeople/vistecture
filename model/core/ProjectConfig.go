package core

import "errors"

type ProjectConfig struct {
	Name           string           `json:"name" yaml:"name" `
	CoreComponents []*CoreComponent `json:"core-components,omitempty" yaml:"core-components,omitempty"`
	Components     []*Component     `json:"components,omitempty" yaml:"components,omitempty"`
}

//Validates Project and Components
func (ProjectConfig *ProjectConfig) Validate() []error {
	var foundErrors []error

	for _, coreComponent := range ProjectConfig.CoreComponents {
		if coreComponent.NoDependency && len(coreComponent.Dependencies) >= 1 {
			foundErrors = append(foundErrors, errors.New("Core component with name " + coreComponent.Name + " has dependencies defined together with NoDependency true"))
		}
	}
	return foundErrors
}

