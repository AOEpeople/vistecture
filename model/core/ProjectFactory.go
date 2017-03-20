package core

import (

	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"strings"
	"path/filepath"

	"errors"
)

// we need this to bind the interface funcs.. Basti?? Best practice?
type ProjectFactory struct {
}


// Factory
func CreateProjectAndValidate(filePath string) (*Project, error) {
	var factory ProjectFactory
	return factory.LoadFromFilePath(filePath)
}

//Loads from JSON file or Folder and returns reference to new Project with all data merged
func (factory *ProjectFactory) LoadFromFilePath(filePath string) (*Project, error) {
	fileStat, e := os.Stat(filePath);
	if  e != nil {
		return nil, errors.New(fmt.Sprintf("No valid filepath - error: %v\n", e))
	}

	switch mode := fileStat.Mode(); {
	case mode.IsDir():
		newProject, e := factory.createFromFolder(filePath)
		if (e != nil) {
			return nil, e
		}
		return newProject,nil
	case mode.IsRegular():
		newProject, e := factory.createFromFile(filePath)
		if (e != nil) {
			return nil, e
		}
		return newProject,nil
	}
	var newProject *Project
	return newProject, nil
}

func (factory *ProjectFactory) createFromFolder(folderPath string) (*Project, error) {
	files, e := filepath.Glob(strings.TrimRight(folderPath,"/")+"/*.json")
	if e != nil {
		return nil, e
	}
	var newProject Project
	if len(files) == 0 {
		return nil, errors.New("No JSON files found in folder \""+folderPath+"\"")
	}
	for _, file := range files {
		projectForFile, e := factory.createFromFile(file)
		if (e != nil) {
			return &newProject, e
		}
		e = newProject.AddComponentsFromProject(projectForFile)
		if e != nil {
			return &newProject, errors.New(e.Error() + " in file "+file)
		}
	}
	return &newProject, nil
}

func (factory *ProjectFactory) createFromFile(fileName string) (*Project, error) {
	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		return nil,errors.New("File error: " + e.Error())
	}
	var newProject Project
	e = json.Unmarshal(file, &newProject)
	if (e != nil) {
		return nil,errors.New(fmt.Sprintf("JSON broken in %v error: %v\n", fileName, e))
	}
	return &newProject, nil
}
