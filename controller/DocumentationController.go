package controller

import (
	"fmt"
	"os"
	"text/template"

	"github.com/danielpoe/appdependency/model/core"
	"github.com/danielpoe/appdependency/model/graphviz"
)

type (
	DocumentationController struct {
		ProjectConfigPath *string
	}

	TemplateData struct {
		Project *core.Project
	}
)

func (DocumentationController DocumentationController) GraphvizAction() {
	Project := loadProject(*DocumentationController.ProjectConfigPath)
	ProjectDrawer := graphviz.CreateProjectDrawer(Project)
	fmt.Print(ProjectDrawer.Draw())
}

func (DocumentationController DocumentationController) HTMLDocumentAction() {
	project := loadProject(*DocumentationController.ProjectConfigPath)
	tpl := template.Must(template.ParseFiles("templates/htmldocument.tmpl"))
	tpl.Execute(os.Stdout, TemplateData{Project: project})
}
