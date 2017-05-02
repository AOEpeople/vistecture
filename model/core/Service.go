package core

type Service struct {
	Name         string       `json:"name" yaml:"name"`
	Description  string       `json:"description" yaml:"description"`
	Type         string       `json:"type,omitempty" yaml:"type,omitempty"`
	IsPublic     bool         `json:"isPublic,omitempty" yaml:"isPublic,omitempty"`
	Dependencies []Dependency `json:"dependencies" yaml:"dependencies"`
}
