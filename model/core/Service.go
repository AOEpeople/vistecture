package core

type Service struct {
	Name          string       `json:"name" yaml:"name"`
	Title         string       `json:"title" yaml:"title"`
	Summary       string       `json:"summary" yaml:"summary"`
	Description   string       `json:"description" yaml:"description"`
	Type          string       `json:"type,omitempty" yaml:"type,omitempty"`
	IsPublic      bool         `json:"isPublic,omitempty" yaml:"isPublic,omitempty"`
	IsOpenHost    bool         `json:"isOpenHost,omitempty" yaml:"isOpenHost,omitempty"`
	SecurityLevel string       `json:"securityLevel" yaml:"securityLevel"`
	Dependencies  []Dependency `json:"dependencies" yaml:"dependencies"`
	Status        string       `json:"status" yaml:"status"`
}
