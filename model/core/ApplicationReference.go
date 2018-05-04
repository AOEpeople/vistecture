package core

type ApplicationReference struct {
	//Name - is used to reference
	Name string `json:"name" yaml:"name"`
	//Title and all other attributes are supposed to override or extend the properties of the referenced application
	Title               string            `json:"title,omitempty" yaml:"title,omitempty"`
	Summary             string            `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description         string            `json:"description,omitempty" yaml:"description,omitempty"`
	Group               string            `json:"group,omitempty" yaml:"group,omitempty"`
	Technology          string            `json:"technology,omitempty" yaml:"technology,omitempty"`
	Category            string            `json:"category,omitempty" yaml:"category,omitempty"`
	AddProvidedServices []Service         `json:"add-provided-services" yaml:"add-provided-services"`
	AddDependencies     []Dependency      `json:"add-dependencies" yaml:"add-dependencies"`
	Properties          map[string]string `json:"properties" yaml:"properties"`
}
