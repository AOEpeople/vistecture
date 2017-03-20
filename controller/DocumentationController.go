package controller

import (
	"fmt"
	"appdependency/model/graphviz"
	"text/template"
	"os"
)

type DocumentationController struct {
	ProjectConfigPath string
}

func (DocumentationController DocumentationController) GraphvizAction() {
	Project:=loadProject(DocumentationController.ProjectConfigPath)
	ProjectDrawer := graphviz.CreateProjectDrawer(Project)
	fmt.Print(ProjectDrawer.Draw())
}


func (DocumentationController DocumentationController) HTMLDocumentAction() {
	Project:=loadProject(DocumentationController.ProjectConfigPath)
	template := template.Must(template.ParseFiles("templates/htmldocument.tmpl"))
	data := map[string]interface{}{
		"projectName": Project.Name,
		"componentTable": Project.AsTable(),
	}
	template.Execute(os.Stdout, data)
}
