package application

import "github.com/AOEpeople/vistecture/model/core"

// Factory for using a dedicated project name
// if projectName is empty it loads the default project
func CreateProjectByName(filePath string, projectName string, skipValidation bool) (*core.Project, []error) {
	var definitionLoader projectDefinitionLoader
	var foundErrors []error
	definitionLoader.strictMode = !skipValidation
	projectDefinitions, error := definitionLoader.LoadFromFilePath(filePath)
	if error != nil {
		foundErrors = append(foundErrors, error)
		return nil, foundErrors
	}

	//Collects errors for repository validation
	foundErrors = append(foundErrors, projectDefinitions.Validate()...)

	project, projectErrors := CreateProjectFromProjectDefinitions(projectDefinitions, projectName)

	//Collects errors for project building
	foundErrors = append(foundErrors, projectErrors...)
	//Collects errors for project validation
	if !skipValidation {
		foundErrors = append(foundErrors, project.Validate()...)
	}
	return project, foundErrors
}

func CreateProjectFromProjectDefinitions(projectDefinitions *ProjectDefinitions, projectName string) (*core.Project, []error) {

	var newProject core.Project

	projectConfig := projectDefinitions.GetProjectConfig(projectName)

	newProject.Name = projectConfig.Name

	//If the project has no explicit defined list of included applications, use all
	if len(projectConfig.IncludedApplication) < 1 {
		newProject.Applications = projectDefinitions.Applications
		newProject.GenerateApplicationIds()
		return &newProject, nil
	}

	//If we have included-applications configured, than apply the applicationReference settings and include only this applications
	var foundErrors []error
	for _, referencedApplication := range projectConfig.IncludedApplication {
		projectApplication, error := projectDefinitions.FindApplicationByName(referencedApplication.Name)
		if error == nil {
			adjustedApplication, err := referencedApplication.GetAdjustedApplication(projectApplication)
			if err != nil {
				foundErrors = append(foundErrors, err)
				continue
			}
			newProject.Applications = append(newProject.Applications, adjustedApplication)
		} else {
			foundErrors = append(foundErrors, error)
		}
	}
	newProject.GenerateApplicationIds()
	return &newProject, foundErrors
}
