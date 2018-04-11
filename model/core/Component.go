package core

import (
	"fmt"
)

type Component struct {
	Name		string		`json:"name" yaml:"name"`
}

func (Component Component) Validate() bool {
	if Component.Name == "" {
		fmt.Printf("Component Name is null")
		return false
	}
	return true
}
