package application

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/AOEpeople/vistecture/v2/model/core"
	yaml "gopkg.in/yaml.v2"
)

type (
	ProjectLoader struct {
		StrictMode bool
	}
	oldApplicationFormat struct {
		Applications []*core.Application `json:"applications" yaml:"applications"`
	}

	ErrorCollection struct {
		Errors []error
	}
)

func (e *ErrorCollection) Error() string {
	return fmt.Sprintf("%v", e.Errors)
}

func (e *ErrorCollection) Add(err error) {
	if err == nil {
		return
	}
	if errMany, ok := err.(*ErrorCollection); ok {
		e.Errors = append(e.Errors, errMany.Errors...)
	} else {
		e.Errors = append(e.Errors, err)
	}
}

func (e *ErrorCollection) ErrorsOrNil() error {
	if len(e.Errors) > 0 {
		return e
	}
	return nil
}

func (p *ProjectLoader) LoadProjectConfig(filePath string) (*ProjectConfig, error) {
	if !strings.Contains(filePath, ".yml") && !!strings.Contains(filePath, ".yaml") {
		return nil, errors.New("wrong fileextension")
	}
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var projectConfig *ProjectConfig
	err = p.unmarshalYaml(file, &projectConfig)

	if err != nil {
		return nil, err
	}
	return projectConfig, nil
}

func (p *ProjectLoader) LoadProjectFromConfigFile(filePath string, limitToSubView string) (*core.Project, error) {
	projectConfig, err := p.LoadProjectConfig(filePath)
	if err != nil {
		return nil, err
	}
	baseFolder := path.Dir(filePath)
	return p.LoadProject(projectConfig, baseFolder, limitToSubView)
}

func (p *ProjectLoader) LoadProject(projectConfig *ProjectConfig, baseFolder string, limitToSubView string) (*core.Project, error) {
	collectedErrors := &ErrorCollection{}
	var newProject core.Project
	newProject.Name = projectConfig.ProjectName

	var applications []*core.Application
	for _, pathsWithAppDefinitions := range projectConfig.AppDefinitionsPaths {
		loadedApps, err := p.LoadApplications(path.Join(baseFolder, pathsWithAppDefinitions))
		if err != nil {
			collectedErrors.Add(err)
		}
		if loadedApps == nil {
			continue
		}
		applications = append(applications, loadedApps...)
	}
	//Apply overrides in config
	for _, override := range projectConfig.AppOverrides {
		app, found := findApplicationByName(override.Name, applications)
		if !found {
			log.Println("Override defined for an unknown application " + override.Name)
			continue
		}

		adjustedApplication, err := override.GetAdjustedApplication(app)
		if err != nil {
			collectedErrors.Add(err)
			continue
		}
		applications = replaceApplication(adjustedApplication, applications)
	}

	if limitToSubView == "" {
		newProject.Applications = applications
		newProject.GenerateApplicationIds()
		return &newProject, collectedErrors.ErrorsOrNil()
	}
	for _, subViewConfig := range projectConfig.SubViewConfig {
		if subViewConfig.Name == limitToSubView {
			newProject.Applications = subViewConfig.GetMatchedApps(applications)
			newProject.GenerateApplicationIds()
			return &newProject, collectedErrors.ErrorsOrNil()
		}
	}
	collectedErrors.Add(errors.New(fmt.Sprintf("Subview with name %v not defined", limitToSubView)))
	return nil, collectedErrors
}

func (p *ProjectLoader) LoadApplications(filePath string) ([]*core.Application, error) {
	collectedErrors := &ErrorCollection{}
	if filePath == "" {
		collectedErrors.Add(errors.New("No applications definitions file path given"))
		return nil, collectedErrors
	}
	var applications []*core.Application
	fileStat, err := os.Stat(filePath)
	if err != nil {
		collectedErrors.Add(errors.New(fmt.Sprintf("No valid filepath (%v) to load application - error: %v\n", filePath, err)))
		return nil, collectedErrors
	}

	switch mode := fileStat.Mode(); {
	case mode.IsDir():
		loadedApplications, err := p.createFromFolder(filePath)
		if err != nil {
			collectedErrors.Add(err)
		}
		if loadedApplications == nil {
			return nil, collectedErrors
		}
		applications = append(applications, loadedApplications...)
	case mode.IsRegular():
		loadedApplications, err := p.createFromFile(filePath)
		if err != nil {
			collectedErrors.Add(err)
		}
		if loadedApplications == nil {
			return nil, collectedErrors
		}
		applications = append(applications, loadedApplications...)
	}
	//Check duplicates
	for sk, sapp := range applications {
		i := 0
		for ck, capp := range applications {
			if sk > ck && sapp.Name == capp.Name {
				i ++
				collectedErrors.Add(fmt.Errorf("Application with name %v is duplicated - exists in %d and %d",sapp.Name,sk,ck))
				applications[ck].Name = fmt.Sprintf("%v-Duplicate-%v",applications[ck].Name,i)
			}
		}
	}

	return applications, collectedErrors
}

func (p *ProjectLoader) createFromFolder(folderPath string) ([]*core.Application, error) {
	collectedErrors := &ErrorCollection{}
	var applications []*core.Application
	files, err := filepath.Glob(strings.TrimRight(folderPath, "/") + "/*")
	if err != nil {
		collectedErrors.Add(err)
		return nil, collectedErrors
	}
	if len(files) == 0 {
		return nil, errors.New("No files found in folder \"" + folderPath + "\"")
	}
	for _, file := range files {

		fileInfo, fileErr := os.Stat(file)
		if fileErr != nil {
			collectedErrors.Add(fileErr)
			continue
		}
		if fileInfo.IsDir() {
			if !strings.Contains(fileInfo.Name(), ".git") {
				loadedApps, err := p.createFromFolder(file)
				if err != nil {
					collectedErrors.Add(err)
				}
				if loadedApps == nil {
					continue
				}
				applications = append(applications, loadedApps...)
			}
		} else if !fileInfo.IsDir() && (strings.Contains(fileInfo.Name(), ".yml") || strings.Contains(fileInfo.Name(), ".yaml")) {
			loadedApps, err := p.createFromFile(file)
			if err != nil {
				collectedErrors.Add(err)
			}
			if loadedApps == nil {
				continue
			}
			applications = append(applications, loadedApps...)
		} else {
			continue
		}
	}
	return applications, collectedErrors.ErrorsOrNil()
}

func (p *ProjectLoader) createFromFile(fileName string) ([]*core.Application, error) {
	var applications []*core.Application
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.New("File error: " + err.Error())
	}
	if !strings.Contains(fileName, ".yml") && !strings.Contains(fileName, ".yaml") {
		return nil, errors.New("Unknown file type")
	}
	//first load yml format where we expect only one application:
	var loadedApplication core.Application
	errNewFormat := p.unmarshalYaml(file, &loadedApplication)
	var oldFormat oldApplicationFormat
	errOldFormat := p.unmarshalYaml(file, &oldFormat)

	//Decide automatically which format should be used to be added to the result:
	if errNewFormat == nil && errOldFormat == nil {
		//both parsing succeeds - use the one with content
		if loadedApplication.Name != "" {
			applications = append(applications, &loadedApplication)
		} else {
			applications = append(applications, oldFormat.Applications...)
		}
	} else if errNewFormat != nil && errOldFormat == nil {
		applications = append(applications, oldFormat.Applications...)
	} else if errNewFormat == nil {
		applications = append(applications, &loadedApplication)
	} else {
		return nil, errors.New(fmt.Sprintf("Cannot parse application definition file %v: \n \t Errors interpreted in 'Single App Format': %v \n \t Errors interpreted in 'Multiple App Format': %v", fileName, errNewFormat, errOldFormat))
	}
	return applications, nil
}

func (p *ProjectLoader) unmarshalYaml(file []byte, i interface{}) error {
	if p.StrictMode {
		return yaml.UnmarshalStrict(file, i)
	} else {
		return yaml.Unmarshal(file, i)
	}
}

func findApplicationByName(name string, apps []*core.Application) (*core.Application, bool) {
	for _, app := range apps {
		if app.Name == name {
			return app, true
		}
	}
	return nil, false
}

func replaceApplication(appToReplace *core.Application, apps []*core.Application) []*core.Application {
	for k, app := range apps {
		if app.Name == appToReplace.Name {
			apps[k] = appToReplace
			break
		}
	}
	return apps
}
