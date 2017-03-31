package core

type Service struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Type         string       `json:"type,omitempty"`
	IsPublic     bool         `json:"isPublic,omitempty"`
	Dependencies []Dependency `json:"dependencies"`
}
