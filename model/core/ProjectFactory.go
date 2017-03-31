package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// we need this to bind the interface funcs.. Basti?? Best practice?
type ProjectFactory struct{}

// Factory
func CreateProject(filePath string) (*Project, error) {
	var factory ProjectFactory
	return factory.LoadFromFilePath(filePath)
}

//Loads from JSON file or Folder and returns reference to new Project with all data merged
func (factory *ProjectFactory) LoadFromFilePath(filePath string) (*Project, error) {
	fileStat, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("No valid filepath - error: %v\n", err))
	}

	var newProject *Project

	switch mode := fileStat.Mode(); {
	case mode.IsDir():
		newProject, err = factory.createFromFolder(filePath)
	case mode.IsRegular():
		newProject, err = factory.createFromFile(filePath)
	}

	if err != nil {
		return nil, err
	}

	return newProject, nil
}

func (factory *ProjectFactory) createFromFolder(folderPath string) (*Project, error) {
	files, err := filepath.Glob(strings.TrimRight(folderPath, "/") + "/*")
	if err != nil {
		return nil, err
	}
	var newProject Project
	if len(files) == 0 {
		return nil, errors.New("No files found in folder \"" + folderPath + "\"")
	}
	for _, file := range files {

		fileInfo, fileErr := os.Stat(file)
		if fileErr != nil {
			return &newProject, fileErr
		}

		var tempProject *Project
		var tmpError error

		if fileInfo.IsDir() {
			tempProject, tmpError = factory.createFromFolder(file)
		} else if !fileInfo.IsDir() && strings.Contains(fileInfo.Name(),".json") {
			tempProject, tmpError = factory.createFromFile(file)
		} else {
			continue
		}

		if tmpError != nil {
			return &newProject, errors.New(tmpError.Error() + " in file " + file)
		}
		err =  newProject.AddComponentsFromProject(tempProject)
		if err != nil {
			return &newProject, errors.New(err.Error() + " in file " + file)
		}
	}
	return &newProject, nil
}

func (factory *ProjectFactory) createFromFile(fileName string) (*Project, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.New("File error: " + err.Error())
	}
	var newProject Project
	err = json.Unmarshal(file, &newProject)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("JSON broken in %v error: %v\n", fileName, err))
	}
	return &newProject, nil
}
