package core

import (
	"strings"
	"fmt"
	"errors"
)

type Component struct {
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	Group string `json:"group,omitempty"`
	Technology string `json:"technology,omitempty"`
	Category string `json:"category,omitempty"`
	ProvidedServices []Service `json:"provided-services"`
	Dependencies []Dependency `json:"dependencies"`
	Display ComponentDisplaySettings `json:"display,omitempty"`
}

//Basti: *Component or Component - PRO CON ?
func (Component Component) Validate() bool {
	if (strings.Contains(Component.Name,".")) {
		fmt.Printf("Component Name contains . '%v'\n", Component.Name)
		return false
	}
	return true
}



func (Component Component) FindService(nameToMatch string) (*Service,error) {
	for _, service := range Component.ProvidedServices {
		if (service.Name == nameToMatch) {
			return &service,nil
		}
	}
	//Basti: is this the way to trigger error?
	return nil,errors.New("Component "+ Component.Name + " has no Interface with Name " +nameToMatch)
}


