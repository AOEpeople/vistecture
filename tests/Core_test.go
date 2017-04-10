package tests

import (
	"testing"
	core "vistecture/model/core"
)


func TestCreateProjectFromFixture(t *testing.T) {

	project, e := core.CreateProject("fixture")
	if e != nil {
		t.Error("Factory returned error",e)
	}
	if project.Name != "Fixture with 2 apps" {
		t.Error("Expected name Fixture with 2 apps")
	}

	application,e  := project.FindApplication("app1")
	if e != nil {
		t.Error("Project returned error when expecting app1",e)
	}
	if (application.Name != "app1") {
		t.Error("Expected application with Name app2")
	}
}



func TestGetReverseDependencies(t *testing.T) {

	project:= core.Project {
		"Project1",
		[]*core.Application {
			{
				"app1",
				"",
				"",
				"",
				"",
				"",
				nil,
				nil,
				[]core.Dependency{
					{
						"app2",
						"",
						false,
						false,
						true,
					},
				},
				core.ApplicationDisplaySettings{},
			},
			{
				"app2",
				"",
				"",
				"",
				"",
				"",
				nil,
				nil,
				nil,
				core.ApplicationDisplaySettings{},
			},
		} ,
	}

	if !contains(project.FindApplicationThatReferenceTo(project.Applications[1], false),project.Applications[0]) {
		t.Error("Expected application1 to link to application2")
	}

	if project.FindApplicationThatReferenceTo(project.Applications[0], false) != nil {
		t.Error("Expected empty slice to reference to application1")
	}
}


func contains(searchIn []*core.Application, findApp *core.Application) bool {
	for _, app := range searchIn {
		if app == findApp {
			return true
		}
	}
	return false
}