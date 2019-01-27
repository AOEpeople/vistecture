package core

import (
	"testing"
)

func TestProject_FindAllApplicationsThatReferenceApplication(t *testing.T) {

	project := Project{
		"Project1",
		[]*Application{
			{
				Name: "app1",

				Dependencies: []Dependency{
					{
						Reference: "app2",
					},
				},
			},
			{
				Name: "app2",
			},
		},
	}

	if !contains(project.FindApplicationThatReferenceTo(project.Applications[1], false), project.Applications[0]) {
		t.Error("Expected application1 to link to application2")
	}

	if project.FindApplicationThatReferenceTo(project.Applications[0], false) != nil {
		t.Error("Expected empty slice to reference to application1")
	}
}

func contains(searchIn []*Application, findApp *Application) bool {
	for _, app := range searchIn {
		if app == findApp {
			return true
		}
	}
	return false
}
