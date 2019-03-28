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
		Name                 string                    `json:"name"`
		AvailableSubViews    []string                  `json:"availableSubViews"`
		ApplicationsByGroup  *core.ApplicationsByGroup `json:"applicationsByGroup"`
		ApplicationsDto      []*ApplicationDto         `json:"applications"`
		StaticDocumentations []string                  `json:"staticDocumentations"`
		Errors               []string                  `json:"errors"`
		MissingApplications  MissingApplications       `json:"missingApplications"`
	}

	ApplicationDto struct {
		*core.Application
		DependenciesGrouped               []*core.DependenciesGrouped `json:"dependenciesGrouped"`
		DependenciesToMissingApplications []*MissingApplicationDto    `json:"dependenciesToMissingApplications"`
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
	subViewName, _ := r.URL.Query()["subview"]
	project, err := p.projectLoader.LoadProject(p.projectDefinitions, p.definitionsBaseFolder, strings.Join(subViewName, ""))
	result := Result{
		Name:                project.Name,
		ApplicationsByGroup: project.GetApplicationsRootGroup(),
	}

	if err != nil {
		result.AddError(err)
	}

	if project == nil {
		result.AddError(errors.New("No project loaded"))
		p.writeJson(w, result, true)
		return
	}

	for _, subViewConfig := range p.projectDefinitions.SubViewConfig {
		result.AvailableSubViews = append(result.AvailableSubViews, subViewConfig.Name)
	}

	allMissingApps := new(MissingApplications)
	for _, app := range project.Applications {

		var dependenciesToMissingApplications []*MissingApplicationDto
		dependenciesToMissingApplications = nil
		for _, missing := range app.GetMissingDependencies(project) {
			var id int
			id = 0
			for _, i := range []byte(missing) {
				id = id + int(i)
			}
			missingApp := &MissingApplicationDto{
				Title:    missing,
				PseudoId: id,
			}
			dependenciesToMissingApplications = append(dependenciesToMissingApplications, missingApp)
			allMissingApps = allMissingApps.Add(missingApp)
		}

		result.ApplicationsDto = append(result.ApplicationsDto, &ApplicationDto{
			Application:                       app,
			DependenciesGrouped:               app.GetDependenciesGrouped(project),
			DependenciesToMissingApplications: dependenciesToMissingApplications,
		})
	}
	files, err := getStaticDocuments(documentsFolder)
	if err != nil {
		result.AddError(err)
	}
	result.MissingApplications = *allMissingApps
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
