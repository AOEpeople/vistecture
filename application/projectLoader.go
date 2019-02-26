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

	"github.com/AOEpeople/vistecture/model/core"
	"gopkg.in/yaml.v2"
)

type (
	ProjectLoader struct {
		StrictMode bool
	}
	oldApplicationFormat struct {
		Applications []*core.Application `json:"applications" yaml:"applications"`
	}
)

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
	if p.StrictMode {
		err = yaml.UnmarshalStrict(file, &projectConfig)
	} else {
		err = yaml.Unmarshal(file, &projectConfig)
	}
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
	var newProject core.Project
	newProject.Name = projectConfig.ProjectName

	var applications []*core.Application
	for _, pathsWithAppDefinitions := range projectConfig.AppDefinitionsPaths {
		loadedApps, err := p.loadApplications(path.Join(baseFolder, pathsWithAppDefinitions))
		if err != nil {
			return nil, err
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
			return nil, err
		}
		applications = replaceApplication(adjustedApplication, applications)
	}

	if limitToSubView == "" {
		newProject.Applications = applications
		newProject.GenerateApplicationIds()
		return &newProject, nil
	}
	for _, subViewConfig := range projectConfig.SubViewConfig {
		if subViewConfig.Name == limitToSubView {
			newProject.Applications = subViewConfig.GetMatchedApps(applications)
			newProject.GenerateApplicationIds()
			return &newProject, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Subview with name %v not defined", limitToSubView))

}

func (p *ProjectLoader) loadApplications(filePath string) ([]*core.Application, error) {
	var applications []*core.Application
	fileStat, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("No valid filepath - error: %v\n", err))
	}

	switch mode := fileStat.Mode(); {
	case mode.IsDir():
		loadedApplications, err := p.createFromFolder(filePath)
		if err != nil {
			return nil, err
		}
		applications = append(applications, loadedApplications...)
	case mode.IsRegular():
		loadedApplications, err := p.createFromFile(filePath)
		if err != nil {
			return nil, err
		}
		applications = append(applications, loadedApplications...)
	}
	return applications, nil
}

func (p *ProjectLoader) createFromFolder(folderPath string) ([]*core.Application, error) {
	var applications []*core.Application
	files, err := filepath.Glob(strings.TrimRight(folderPath, "/") + "/*")
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.New("No files found in folder \"" + folderPath + "\"")
	}
	for _, file := range files {

		fileInfo, fileErr := os.Stat(file)
		if fileErr != nil {
			return nil, fileErr
		}
		if fileInfo.IsDir() {
			if !strings.Contains(fileInfo.Name(), ".git") {
				loadedApps, err := p.createFromFolder(file)
				if err != nil {
					return nil, err
				}
				applications = append(applications, loadedApps...)
			}
		} else if !fileInfo.IsDir() && (strings.Contains(fileInfo.Name(), ".yml") || strings.Contains(fileInfo.Name(), ".yaml")) {
			loadedApps, err := p.createFromFile(file)
			if err != nil {
				return nil, err
			}
			applications = append(applications, loadedApps...)
		} else {
			continue
		}
	}
	return applications, nil
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
			log.Println("You use DEPRICATED Old format for file " + fileName)
			applications = append(applications, oldFormat.Applications...)
		}
	} else if errNewFormat != nil && errOldFormat == nil {
		log.Println("You use DEPRICATED Old format for file " + fileName)
		applications = append(applications, oldFormat.Applications...)
	} else if errNewFormat == nil {
		applications = append(applications, &loadedApplication)
	} else {
		return nil, errors.New(fmt.Sprintf("Cannot parse file %v / %v", errNewFormat, errOldFormat))
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

func replaceApplication(app *core.Application, apps []*core.Application) []*core.Application {
	for k, app := range apps {
		if app.Name == app.Name {
			apps[k] = app
		}
	}
	return apps
}
