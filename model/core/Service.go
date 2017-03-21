package core

type Service struct {
	Name         string       `json:"name"`
	Type         string       `json:"type,omitempty"`
	IsPublic     bool         `json:"isPublic,omitempty"`
	Dependencies []Dependency `json:"dependencies"`
}
