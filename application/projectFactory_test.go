package application

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateProjectByName(t *testing.T) {
	project, errors := CreateProjectByName("fixtures/fixture-json", "", false)
	if len(errors) >= 1 {
		t.Fatalf("Factory returned error: %s", errors)
	}
	if project.Name != "Fixture Project" {
		t.Error("Expected name 'Fixture Project' Got " + project.Name)
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

	project, errors := CreateProjectByName("fixtures/fixture-merge", "", false)
	if len(errors) >= 1 {
		t.Fatalf("Factory returned error: %s", errors)
	}
	if project.Name != "Fixture Project Merge" {
		t.Error("Expected name 'Fixture Project Merge'")
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

func TestCreateProjectFromMultiple1(t *testing.T) {

	project, errors := CreateProjectByName("fixtures/fixture-multiple", "Fixture project Multiple 1", false)
	if len(errors) >= 1 {
		t.Fatalf("Factory returned error: %s", errors)
	}
	if project.Name != "Fixture Project Multiple 1" {
		t.Error("Expected name 'Fixture Project Multiple 1'")
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

	assert.Equal(t, "core", application3.Category)

}

func TestCreateProjectFromMultiple2(t *testing.T) {

	project, errors := CreateProjectByName("fixtures/fixture-multiple", "Fixture Project Multiple 2", false)
	if len(errors) >= 1 {
		t.Fatalf("Factory returned error: %s", errors)
	}
	if project.Name != "Fixture Project Multiple 2" {
		t.Error("Expected name 'Fixture project Multiple 2' / Got" + project.Name)
	}

	application, e := project.FindApplication("app4")
	if e == nil {
		t.Error("Expected application app4 to be missing but is available:" + application.Name)
	}
}

func TestCreateProjectFromBoProject(t *testing.T) {

	project, errors := CreateProjectByName("fixtures/fixture-noproject", "", false)
	if len(errors) >= 1 {
		t.Fatalf("Factory returned error: %s", errors)
	}
	assert.Equal(t, "Full Project Definitions", project.Name)

	application, e := project.FindApplication("app5")
	if e != nil {
		t.Error("project returned error when expecting app5", e)
	}
	assert.Equal(t, "app5", application.Name)
}

func TestNoDefinitionFound(t *testing.T) {

	_, errors := CreateProjectByName("fake-dir", "", false)
	if errors == nil {
		t.Error("Expected errors to be filled", errors)
	}
	if strings.Contains(errors[0].Error(), "Could not build repository: No files found in folder") {
		t.Error("Expected error: 'Could not build repository: No files found in folder'")
	}
}

func TestExampleProjects(t *testing.T) {

	project, errors := CreateProjectByName("../example/demoproject", "", false)
	if len(errors) >= 1 {
		t.Fatalf("Factory returned error: %s", errors)
	}
	if project.Name != "Demoproject" {
		t.Error("Expected name 'Demoproject'")
	}

}
