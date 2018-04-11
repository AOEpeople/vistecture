package core

// we need this to bind the interface funcs
type ProjectFactory struct{}

// Factory
func CreateProject(filePath string) (*Project, []error) {
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
	foundErrors = append(foundErrors, project.Validate()...)

	return project, foundErrors
}

// Factory for using a dedicated project name
func CreateProjectByName(filePath string, projectName string) (*Project, []error) {
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
	foundErrors = append(foundErrors, project.Validate()...)

	return project, foundErrors
}

func CreateProjectFromRepository(repository *Repository, projectName string) (*Project, []error) {

	var newProject Project
	var foundErrors []error


	projectInfo := repository.GetProjectInfo(projectName)

	newProject.Name = projectInfo.Name

	//Filter project components if defined. If not all repo applications will be used.
	if len(projectInfo.Components) >= 1 {
		for _, component := range projectInfo.Components {
			application, error := repository.FindApplicationByName(component.Name)
			if error == nil {
				newProject.Applications = append(newProject.Applications, &application)
			} else {
				foundErrors = append(foundErrors, error)
			}
		}
	} else {
		components := repository.FindNonCoreApplications()
		newProject.Applications = components
	}


	//Filter project for core components if any is defined
	if len(projectInfo.CoreComponents) >= 1 {
		for _, component := range projectInfo.CoreComponents {
			application, error := repository.FindApplicationByName(component.Name)
			if error == nil {
				dependencies := component.Dependencies
				if component.NoDependency || len(dependencies) > 0 {
					application.Dependencies = component.Dependencies
				}
				newProject.Applications = append(newProject.Applications, &application)
			} else {
				foundErrors = append(foundErrors, error)
			}
		}
	} else {
		for _, component := range repository.FindApplicationsByCategory(CORE) {
			newProject.Applications = append(newProject.Applications, component)
		}
	}
	return &newProject, foundErrors
}