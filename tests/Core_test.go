package tests

import (
	"testing"

	core "github.com/AOEpeople/vistecture/model/core"
)

func TestCreateProjectFromFixture(t *testing.T) {

	project, e := core.CreateProject("fixture")
	if e != nil {
		t.Error("Factory returned error", e)
	}
	if project.Name != "Fixture with 2 apps" {
		t.Error("Expected name Fixture with 2 apps")
	}

	application, e := project.FindApplication("app1")
	if e != nil {
		t.Error("Project returned error when expecting app1", e)
	}
	if application.Name != "app1" {
		t.Error("Expected application with Name app2")
	}
}

func TestCreateProjectFromFixtureFolderWithMerge(t *testing.T) {

	project, e := core.CreateProject("fixture-merge")
	if e != nil {
		t.Error("Factory returned error", e)
	}
	if project.Name != "test2" {
		t.Error("Expected name test2")
	}

	application, e := project.FindApplication("app1")
	if e != nil {
		t.Error("Project returned error when expecting app1", e)
	}
	if application.Name != "app1" {
		t.Error("Expected application with Name app2")
	}

	if application.Properties["git"] != "here" {
		t.Error("Expected property git with value here")
	}
}

func TestGetReverseDependencies(t *testing.T) {

	project := core.Project{
		"Project1",
		[]*core.Application{
			{
				Name: "app1",

				Dependencies: []core.Dependency{
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

func contains(searchIn []*core.Application, findApp *core.Application) bool {
	for _, app := range searchIn {
		if app == findApp {
			return true
		}
	}
	return false
}
