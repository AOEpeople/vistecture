package application_test

import (
	"testing"

	"github.com/AOEpeople/vistecture/v2/application"
)

func TestProjectLoader_LoadProjectFromConfigFile(t *testing.T) {

	loader := application.ProjectLoader{StrictMode: true}
	project, err := loader.LoadProjectFromConfigFile("fixtures/project.yml", "")
	if err != nil {
		t.Fatal(err)
	}

	_, err = project.FindApplication("app1")
	if err != nil {
		t.Error("expected no error for getting app1 got " + err.Error())
	}

	_, err = project.FindApplication("app2")
	if err != nil {
		t.Error("expected no error for getting app2 got " + err.Error())
	}

}
