package tests

import (
	"testing"

	"strings"

	core "github.com/AOEpeople/vistecture/model/core"
)

func TestCreateProjectFromFixture(t *testing.T) {

	project, errors := core.CreateProject("fixture")
	if len(errors) >= 1 {
		t.Error("Factory returned error", errors)
	}
	if project.Name != "Fixture project" {
		t.Error("Expected name 'Fixture project'")
	}

	application, e := project.FindApplication("app1")
	if e != nil {
		t.Error("project returned error when expecting app1", e)
	}
	if application.Name != "app1" {
		t.Error("Expected application with Name app2")
	}
}

func TestCreateProjectFromFixtureFolderWithMerge(t *testing.T) {

	project, errors := core.CreateProject("fixture-merge")
	if len(errors) >= 1 {
		t.Error("Factory returned error", errors)
	}
	if project.Name != "Fixture project Merge" {
		t.Error("Expected name 'Fixture project Merge'")
	}

	application, e := project.FindApplication("app1")
	if e != nil {
		t.Error("project returned error when expecting app1", e)
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

func TestCreateProjectFromMultiple1(t *testing.T) {

	project, errors := core.CreateProjectByName("fixture-multiple", "Fixture project Multiple 1")
	if len(errors) >= 1 {
		t.Error("Factory returned error", errors)
	}
	if project.Name != "Fixture project Multiple 1" {
		t.Error("Expected name 'Fixture project Multiple 1'")
	}

	application, e := project.FindApplication("app1")
	if e != nil {
		t.Error("project returned error when expecting app1", e)
	}
	if application.Name != "app1" {
		t.Error("Expected application with Name app1")
	}

	if application.Properties["git"] != "here" {
		t.Error("Expected property git with value here")
	}

	application3, e := project.FindApplication("app3")
	if e != nil {
		t.Error("project returned error when expecting app3", e)
	}

	if application3.Category != core.CORE.Value() {
		t.Error("Expected category for app3 was core", e)
	}
}

func TestCreateProjectFromMultiple2(t *testing.T) {

	project, errors := core.CreateProjectByName("fixture-multiple", "Fixture project Multiple 2")
	if len(errors) >= 1 {
		t.Error("Factory returned error", errors)
	}
	if project.Name != "Fixture project Multiple 2" {
		t.Error("Expected name 'Fixture project Multiple 2'")
	}

	application, e := project.FindApplication("app4")
	if e == nil {
		t.Error("Expected application app4 to be missing but is available:" + application.Name)
	}
}

func TestCreateProjectFromBoProject(t *testing.T) {

	project, errors := core.CreateProject("fixture-noproject")
	if len(errors) >= 1 {
		t.Error("Factory returned error", errors)
	}
	if project.Name != "Full Repository" {
		t.Error("Expected name 'Full Repository'")
	}

	application, e := project.FindApplication("app5")
	if e != nil {
		t.Error("project returned error when expecting app5", e)
	}
	if application.Name != "app5" {
		t.Error("Expected application with Name app5")
	}
}

func TestCreateReadmeExample(t *testing.T) {

	_, errors := core.CreateProjectByName("fixture-readme", "Ports and Adapters DDD Architecture")
	if len(errors) >= 1 {
		t.Error("Factory returned error", errors)
	}
	project, errors := core.CreateProjectByName("fixture-readme", "Ports and Adapters DDD Architecture minimum")
	if len(errors) >= 1 {
		t.Error("Factory returned error", errors)
	}

	application, _ := project.FindApplication("domain")
	if application.Technology != "play" {
		t.Error("Expected applications technology to be the value 'play', but was: " + application.Technology)
	}

}

func TestNoDefinitionFound(t *testing.T) {

	_, errors := core.CreateProjectByName("fake-dir", "")
	if errors == nil {
		t.Error("Expected errors to be filled", errors)
	}
	if strings.Contains(errors[0].Error(), "Could not build repository: No files found in folder") {
		t.Error("Expected error: 'Could not build repository: No files found in folder'")
	}
}

func TestExampleProjects(t *testing.T) {

	project, errors := core.CreateProject("../example/demoproject")
	if len(errors) >= 1 {
		t.Error("Factory returned error", errors)
	}
	if project.Name != "Demo project" {
		t.Error("Expected name 'Demo project'")
	}

	project, errors = core.CreateProject("../example/ports-and-adapters-architecture")
	if len(errors) >= 1 {
		t.Error("Factory returned error", errors)
	}
	if project.Name != "Ports and Adapters DDD Architecture" {
		t.Error("Expected name 'Ports and Adapters DDD Architecture'")
	}

	project, errors = core.CreateProject("../example/sample-architecture")
	if len(errors) >= 1 {
		t.Error("Factory returned error", errors)
	}
	if project.Name != "Sample" {
		t.Error("Expected name 'Sample'")
	}

}
