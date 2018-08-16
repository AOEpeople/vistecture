package core

// we need this to bind the interface funcs
type ProjectFactory struct{}

// Factory
func CreateProject(filePath string, skipValidation bool) (*Project, []error) {
	var factory RepositoryFactory
	var foundErrors []error
	var projectName string

	repository, error := factory.LoadFromFilePath(filePath)
	if error != nil {
		foundErrors = append(foundErrors, error)
		return &Project{}, foundErrors
	}

	//Collects errors for repository validation
	foundErrors = append(foundErrors, repository.Validate()...)

	project, projectErrors := CreateProjectFromRepository(repository, projectName)
	//Collects errors for project building
	foundErrors = append(foundErrors, projectErrors...)
	//Collects errors for project validation
	if !skipValidation {
		foundErrors = append(foundErrors, project.Validate()...)
	}

	return project, foundErrors
}

// Factory for using a dedicated project name
func CreateProjectByName(filePath string, projectName string, skipValidation bool) (*Project, []error) {
	var factory RepositoryFactory
	var foundErrors []error

	repository, error := factory.LoadFromFilePath(filePath)
	if error != nil {
		foundErrors = append(foundErrors, error)
		return &Project{}, foundErrors
	}

	//Collects errors for repository validation
	foundErrors = append(foundErrors, repository.Validate()...)

	project, projectErrors := CreateProjectFromRepository(repository, projectName)
	//Collects errors for project building
	foundErrors = append(foundErrors, projectErrors...)
	//Collects errors for project validation
	if !skipValidation {
		foundErrors = append(foundErrors, project.Validate()...)
	}
	return project, foundErrors
}

func CreateProjectFromRepository(repository *Repository, projectName string) (*Project, []error) {

	var newProject Project
	var foundErrors []error

	projectConfig := repository.GetProjectConfig(projectName)

	newProject.Name = projectConfig.Name

	if len(projectConfig.IncludedApplication) >= 1 {
		for _, referencedApplication := range projectConfig.IncludedApplication {
			projectApplication, error := repository.FindApplicationByName(referencedApplication.Name)
			if error == nil {
				mergedApplication, err := projectApplication.GetMerged(*referencedApplication)
				if err != nil {
					foundErrors = append(foundErrors, err)
					continue
				}
				newProject.Applications = append(newProject.Applications, &mergedApplication)
			} else {
				foundErrors = append(foundErrors, error)
			}
		}
	} else {
		newProject.Applications = repository.Applications
	}
	return &newProject, foundErrors
}
