package core

import (
	"html/template"
	"strings"

	"github.com/russross/blackfriday"
)

type Dependency struct {
	Reference      string            `json:"reference" yaml:"reference"`
	Description    string            `json:"description" yaml:"description"`
	Relationship   string            `json:"relationship" yaml:"relationship"`
	IsSameLevel    bool              `json:"isSameLevel" yaml:"isSameLevel"`
	IsBrowserBased bool              `json:"isBrowserBased" yaml:"isBrowserBased"`
	Status         string            `json:"status" yaml:"status"`
	Properties     map[string]string `json:"properties" yaml:"properties"`
	IsOptional     bool              `json:"isOptional" yaml:"isOptional"`
}

// Returns the name of the "component" and "service" this dependecy points to
// service might be empty if the dependency just defined the component
func (Dependency *Dependency) GetApplicationAndServiceNames() (string, string) {
	if strings.Contains(Dependency.Reference, ".") {
		splitted := strings.Split(Dependency.Reference, ".")
		return splitted[0], splitted[1]
	}
	return Dependency.Reference, ""
}

func (Dependency *Dependency) GetApplicationName() string {
	if strings.Contains(Dependency.Reference, ".") {
		splitted := strings.Split(Dependency.Reference, ".")
		return splitted[0]
	}
	return Dependency.Reference
}

func (Dependency *Dependency) GetServiceName() string {
	_, s := Dependency.GetApplicationAndServiceNames()
	return s
}

func (Dependency *Dependency) GetApplication(Project *Project) (*Application, error) {
	componentName, _ := Dependency.GetApplicationAndServiceNames()
	return Project.FindApplication(componentName)
}

//GetDescriptionHtml - helper that renders the description text as markdown - to be used in HTML documentations
func (Dependency *Dependency) GetDescriptionHtml() template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(Dependency.Description)))
}
