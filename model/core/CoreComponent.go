package core

import (
	"fmt"
)

type CoreComponent struct {
	Name			string			`json:"name" yaml:"name"`
	NoDependency		bool			`json:"no-dependency" yaml:"no-dependency"`
	Dependencies	[]Dependency    `json:"dependencies" yaml:"dependencies"`
}

func (CoreComponent CoreComponent) Validate() bool {
	if CoreComponent.Name == "" {
		fmt.Printf("CoreComponent Name is null")
		return false
	}
	return true
}
