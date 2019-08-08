package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"fmt"

	"strings"

	"log"

	"path/filepath"

	"github.com/AOEpeople/vistecture/v2/application"
	"github.com/AOEpeople/vistecture/v2/model/core"
	"github.com/gobuffalo/packr/v2"
)

type (
	ProjectController struct {
		projectDefinitions    *application.ProjectConfig
		projectLoader         *application.ProjectLoader
		definitionsBaseFolder string
		skipValidation        bool
	}

	Result struct {
		Name                   string                    `json:"name"`
		AvailableSubViews      []string                  `json:"availableSubViews"`
		ApplicationsByGroup    *core.ApplicationsByGroup `json:"applicationsByGroup"`
		AvailableGroups        *AvailableGroups          `json:"availableGroups"`
		ApplicationsDto        []*ApplicationDto         `json:"applications"`
		StaticDocumentations   []string                  `json:"staticDocumentations"`
		Errors                 []string                  `json:"errors"`
		//MissingApplications - list of applications that are referenced but not definied at all in the projecr
		MissingApplications    MissingApplications       `json:"missingApplications"`
		//UnincludedApplications - list of applications that are referenced but not included in current selection (e.g. because of selected subview or due to a filter)
		UnincludedApplications MissingApplications       `json:"unincludedApplications"`
	}

	AvailableGroups struct {
		SubGroups          []*AvailableGroups `json:"subGroups"`
		GroupName          string             `json:"groupName"`
		QualifiedGroupName string             `json:"qualifiedGroupName"`
	}

	ApplicationDto struct {
		*core.Application
		DependenciesGrouped                  []*core.DependenciesGrouped `json:"dependenciesGrouped"`
		DependenciesToMissingApplications    []*MissingApplicationDto    `json:"dependenciesToMissingApplications"`
		DependenciesToUnincludedApplications []*MissingApplicationDto    `json:"dependenciesToUnincludedApplications"`
	}

	MissingApplications []*MissingApplicationDto

	MissingApplicationDto struct {
		//PseudoId used as id for frontend
		PseudoId int    `json:"id"`
		Title    string `json:"name"`
	}
)

var (
	fileServerInstance http.Handler
)

func (p *ProjectController) Inject(definitions *application.ProjectConfig, projectLoader *application.ProjectLoader, definitionsBaseFolder string, skipValidation bool) {
	p.projectDefinitions = definitions
	p.projectLoader = projectLoader
	p.definitionsBaseFolder = definitionsBaseFolder
	p.skipValidation = skipValidation
}

func (p *ProjectController) IndexAction(w http.ResponseWriter, r *http.Request, localTemplateFolder string) {

	handler := initFileServerInstance(localTemplateFolder)
	handler.ServeHTTP(w, r)
}

func initFileServerInstance(localFolder string) http.Handler {
	if fileServerInstance != nil {
		return fileServerInstance
	}
	var fileSystem http.FileSystem
	if localFolder != "" {
		log.Printf("Using filesystem %v templates for serving", localFolder)
		if _, err := os.Stat(localFolder); os.IsNotExist(err) {
			panic(fmt.Sprintf("Cannot start - Folder %v not exitend", localFolder))
		}
		fileSystem = http.Dir(localFolder)
	} else {
		log.Printf("Using templateBox templates for serving")
		fileSystem = packr.New("templateBox", "./template")
	}
	fileServerInstance = http.FileServer(fileSystem)
	return fileServerInstance
}

func (p *ProjectController) DataAction(w http.ResponseWriter, r *http.Request, documentsFolder string) {
	result := Result{}

	subViewName, _ := r.URL.Query()["subview"]
	completeProject, err := p.projectLoader.LoadProject(p.projectDefinitions, p.definitionsBaseFolder, "")
	if err != nil {
		result.AddError(err)
	}
	if completeProject == nil {
		p.writeJson(w, result, false)
		return
	}
	project, err := p.projectLoader.LoadProject(p.projectDefinitions, p.definitionsBaseFolder, strings.Join(subViewName, ""))
	if err != nil {
		result.AddError(err)
	}
	if project == nil {
		p.writeJson(w, result, false)
		return
	}
	result.AvailableGroups = getAvailableGroups(project.GetApplicationsRootGroup())
	//Filter by filterGroups if parameter is given:
	filterGroupsParam, _ := r.URL.Query()["filterGroups"]
	allGroupFilters := strings.Join(filterGroupsParam, ",")
	if allGroupFilters != "" {
		filterGroups := strings.Split(allGroupFilters, ",")
		var filteredApplications []*core.Application
		for _, app := range project.Applications {
			if inSlice(app.Group, filterGroups) {
				filteredApplications = append(filteredApplications, app)
			}
		}
		project = &core.Project{
			Name:         project.Name,
			Applications: filteredApplications,
		}
	}

	result.Name = project.Name
	result.ApplicationsByGroup = project.GetApplicationsRootGroup()

	if project == nil {
		result.AddError(errors.New("No project loaded"))
		p.writeJson(w, result, true)
		return
	}

	for _, subViewConfig := range p.projectDefinitions.SubViewConfig {
		result.AvailableSubViews = append(result.AvailableSubViews, subViewConfig.Name)
	}

	allMissingApps := new(MissingApplications)
	allUnincludedApps := new(MissingApplications)
	for _, app := range project.Applications {
		var dependenciesToMissingApplications []*MissingApplicationDto
		var dependenciesToUnincludedApplications []*MissingApplicationDto

		//anonymous helper func to create missing app struct
		newMissingApp := func(title string) *MissingApplicationDto {
			var id int
			for _, i := range []byte(title) {
				id = id + int(i)
			}
			return &MissingApplicationDto{
				Title:    title,
				PseudoId: id,
			}
		}

		//anonymous helper func check if app is in given list
		isInList := func(searchFor string, searchIn []*MissingApplicationDto) bool {
			for _, m := range searchIn {
				if m.Title == searchFor {
					return true
				}
			}
			return false
		}

		for _, missing := range app.GetMissingDependencies(completeProject) {
			dependenciesToMissingApplications = append(dependenciesToMissingApplications, newMissingApp(missing))
			allMissingApps = allMissingApps.Add(newMissingApp(missing))
		}
		for _, missing := range app.GetMissingDependencies(project) {
			if isInList(missing, dependenciesToMissingApplications) {
				continue
			}
			dependenciesToUnincludedApplications = append(dependenciesToUnincludedApplications, newMissingApp(missing))
			allUnincludedApps = allUnincludedApps.Add(newMissingApp(missing))
		}

		result.ApplicationsDto = append(result.ApplicationsDto, &ApplicationDto{
			Application:                          app,
			DependenciesGrouped:                  app.GetDependenciesGrouped(project),
			DependenciesToMissingApplications:    dependenciesToMissingApplications,
			DependenciesToUnincludedApplications: dependenciesToUnincludedApplications,
		})
	}
	files, err := getStaticDocuments(documentsFolder)
	if err != nil {
		result.AddError(err)
	}
	result.MissingApplications = *allMissingApps
	result.UnincludedApplications = *allUnincludedApps
	result.StaticDocumentations = files

	p.writeJson(w, result, false)

}

func (p *ProjectController) writeJson(w http.ResponseWriter, result Result, isHardError bool) {
	b, err := json.Marshal(result)
	if err != nil {
		fmt.Fprint(w, "unexpected error: "+err.Error())
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:63343")
	w.Header().Set("Access-Control-Allow-Origin", "null")
	if isHardError {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	fmt.Fprint(w, string(b))
}

func getStaticDocuments(folder string) ([]string, error) {
	var result []string
	if folder == "" {
		return result, nil
	}

	files, err := filepath.Glob(strings.TrimRight(folder, "/") + "/*")
	if err != nil {
		return nil, err
	}
	for _, file := range files {

		fileInfo, fileErr := os.Stat(file)
		if fileErr != nil {
			return nil, fileErr
		}

		if fileInfo.IsDir() {
			continue
		}
		path := strings.TrimPrefix(file, folder)
		result = append(result, path)

	}

	return result, err
}

func (r *Result) AddError(err error) {
	if errorCollection, ok := err.(*application.ErrorCollection); ok {
		for _, singleErr := range errorCollection.Errors {
			r.Errors = append(r.Errors, singleErr.Error())
		}
	} else {
		r.Errors = append(r.Errors, err.Error())
	}

}

func (m MissingApplications) Add(dto *MissingApplicationDto) *MissingApplications {
	for _, ma := range m {
		if ma.PseudoId == dto.PseudoId {
			return &m
		}
	}
	m = append(m, dto)
	return &m
}

func inSlice(search string, in []string) bool {
	for _, v := range in {
		if v == search {
			return true
		}
	}
	return false
}

func getAvailableGroups(group *core.ApplicationsByGroup) *AvailableGroups {
	ag := &AvailableGroups{
		QualifiedGroupName: group.QualifiedGroupName,
		GroupName:          group.GroupName,
	}
	for _, subG := range group.SubGroups {
		ag.SubGroups = append(ag.SubGroups, getAvailableGroups(subG))
	}
	return ag
}
