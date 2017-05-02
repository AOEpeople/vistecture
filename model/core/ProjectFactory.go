package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// we need this to bind the interface funcs
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
		} else if !fileInfo.IsDir() && (strings.Contains(fileInfo.Name(), ".json") || strings.Contains(fileInfo.Name(), ".yml")) {
			tempProject, tmpError = factory.createFromFile(file)
		} else {
			continue
		}

		if tmpError != nil {
			return &newProject, errors.New(tmpError.Error() + " in file " + file)
		}
		err = newProject.MergeWith(tempProject)
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
	if strings.Contains(fileName, ".json") {
		err = json.Unmarshal(file, &newProject)
	} else if strings.Contains(fileName, ".yml") {
		err = yaml.Unmarshal(file, &newProject)
	} else {
		err = errors.New("Unknown file type")
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("File broken in %v error: %v\n", fileName, err))
	}
	return &newProject, nil
}
