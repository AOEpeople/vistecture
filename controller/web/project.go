package web

import (
	"encoding/json"
	"net/http"

	"fmt"

	"strings"

	"github.com/AOEpeople/vistecture/application"
	"github.com/AOEpeople/vistecture/model/core"
)

type (
	ProjectController struct {
		projectDefinitions *application.ProjectDefinitions
	}

	Result struct {
		Name                  string                    `json:"name"`
		AvailableProjectNames []string                  `json:"availableProjectNames"`
		ApplicationsByGroup   *core.ApplicationsByGroup `json:"applicationsByGroup"`
		ApplicationsDto       []*ApplicationDto         `json:"applications"`
	}

	ApplicationDto struct {
		*core.Application
		DependenciesGrouped []*core.DependenciesGrouped `json:"dependenciesGrouped"`
	}
)

func (p *ProjectController) Inject(definitions *application.ProjectDefinitions) {
	p.projectDefinitions = definitions
}

func (p *ProjectController) ShowProjectAction() {

}

func (p *ProjectController) IndexAction(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("templates/web"))
	fs.ServeHTTP(w, r)
}

func (p *ProjectController) DataAction(w http.ResponseWriter, r *http.Request) {
	projectName, _ := r.URL.Query()["project"]
	project, errors := application.CreateProjectFromProjectDefinitions(p.projectDefinitions, strings.Join(projectName, ""))

	if errors != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"Error":"%v"}`, errors)
	}
	result := Result{
		Name:                project.Name,
		ApplicationsByGroup: project.GetApplicationsRootGroup(),
	}
	for _, pConfig := range p.projectDefinitions.ProjectConfig {
		result.AvailableProjectNames = append(result.AvailableProjectNames, pConfig.Name)
	}

	for _, app := range project.Applications {
		result.ApplicationsDto = append(result.ApplicationsDto, &ApplicationDto{
			Application:         app,
			DependenciesGrouped: app.GetDependenciesGrouped(project),
		})
	}
	b, err := json.Marshal(result)
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:63343")
	w.Header().Set("Access-Control-Allow-Origin", "null")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"Error":"`+err.Error()+`"}`)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))

}
