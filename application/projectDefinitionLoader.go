package application

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
type projectDefinitionLoader struct {
	strictMode bool
}

// Main Load/Factory Method
func LoadProjectDefinitions(filePath string) (*ProjectDefinitions, error) {
	var factory projectDefinitionLoader
	return factory.LoadFromFilePath(filePath)
}

//Loads from JSON file or Folder and returns reference to new project with all data merged
func (factory *projectDefinitionLoader) LoadFromFilePath(filePath string) (*ProjectDefinitions, error) {
	fileStat, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("No valid filepath - error: %v\n", err))
	}

	var newRepository *ProjectDefinitions

	switch mode := fileStat.Mode(); {
	case mode.IsDir():
		newRepository, err = factory.createFromFolder(filePath)
	case mode.IsRegular():
		newRepository, err = factory.createFromFile(filePath)
	}

	if err != nil {
		return nil, err
	}

	return newRepository, nil
}

func (factory *projectDefinitionLoader) createFromFolder(folderPath string) (*ProjectDefinitions, error) {
	files, err := filepath.Glob(strings.TrimRight(folderPath, "/") + "/*")
	if err != nil {
		return nil, err
	}
	var newRepository ProjectDefinitions
	if len(files) == 0 {
		return nil, errors.New("No files found in folder \"" + folderPath + "\"")
	}
	for _, file := range files {

		fileInfo, fileErr := os.Stat(file)
		if fileErr != nil {
			return &newRepository, fileErr
		}

		var tempRepository *ProjectDefinitions
		var tmpError error

		if fileInfo.IsDir() {
			if !strings.Contains(fileInfo.Name(), ".git") {
				tempRepository, tmpError = factory.createFromFolder(file)
			}
		} else if !fileInfo.IsDir() && (strings.Contains(fileInfo.Name(), ".json") || (strings.Contains(fileInfo.Name(), ".yml") || strings.Contains(fileInfo.Name(), ".yml"))) {
			tempRepository, tmpError = factory.createFromFile(file)
		} else {
			continue
		}

		if tmpError != nil {
			return &newRepository, errors.New(tmpError.Error() + " in file " + file)
		}
		err = newRepository.mergeWith(tempRepository)
		if err != nil {
			return &newRepository, errors.New(err.Error() + " in file " + file)
		}
	}
	return &newRepository, nil
}

func (factory *projectDefinitionLoader) createFromFile(fileName string) (*ProjectDefinitions, error) {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, errors.New("File error: " + err.Error())
	}
	var newRepository ProjectDefinitions
	if strings.Contains(fileName, ".json") {
		err = json.Unmarshal(file, &newRepository)
	} else if strings.Contains(fileName, ".yml") || strings.Contains(fileName, ".yaml") {
		if factory.strictMode {
			err = yaml.UnmarshalStrict(file, &newRepository)
		} else {
			err = yaml.Unmarshal(file, &newRepository)
		}
	} else {
		err = errors.New("Unknown file type")
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("File broken in %v error: %v\n", fileName, err))
	}
	return &newRepository, nil
}
