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
type RepositoryFactory struct{}

// Factory
func CreateRepository(filePath string) (*Repository, error) {
	var factory RepositoryFactory
	return factory.LoadFromFilePath(filePath)
}

//Loads from JSON file or Folder and returns reference to new project with all data merged
func (factory *RepositoryFactory) LoadFromFilePath(filePath string) (*Repository, error) {
	fileStat, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("No valid filepath - error: %v\n", err))
	}

	var newRepository *Repository

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

func (factory *RepositoryFactory) createFromFolder(folderPath string) (*Repository, error) {
	files, err := filepath.Glob(strings.TrimRight(folderPath, "/") + "/*")
	if err != nil {
		return nil, err
	}
	var newRepository Repository
	if len(files) == 0 {
		return nil, errors.New("No files found in folder \"" + folderPath + "\"")
	}
	for _, file := range files {

		fileInfo, fileErr := os.Stat(file)
		if fileErr != nil {
			return &newRepository, fileErr
		}

		var tempRepository *Repository
		var tmpError error

		if fileInfo.IsDir() {
			//TODO: Configurable ignore paths
			if !strings.Contains(fileInfo.Name(), ".git") {
				tempRepository, tmpError = factory.createFromFolder(file)
			}
		} else if !fileInfo.IsDir() && (strings.Contains(fileInfo.Name(), ".json") || strings.Contains(fileInfo.Name(), ".yml")) {
			tempRepository, tmpError = factory.createFromFile(file)
		} else {
			continue
		}

		if tmpError != nil {
			return &newRepository, errors.New(tmpError.Error() + " in file " + file)
		}
		err = newRepository.MergeWith(tempRepository)
		if err != nil {
			return &newRepository, errors.New(err.Error() + " in file " + file)
		}
	}
	return &newRepository, nil
}

func (factory *RepositoryFactory) createFromFile(fileName string) (*Repository, error) {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, errors.New("File error: " + err.Error())
	}
	var newRepository Repository
	if strings.Contains(fileName, ".json") {
		err = json.Unmarshal(file, &newRepository)
	} else if strings.Contains(fileName, ".yml") {
		err = yaml.Unmarshal(file, &newRepository)
	} else {
		err = errors.New("Unknown file type")
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("File broken in %v error: %v\n", fileName, err))
	}
	return &newRepository, nil
}
