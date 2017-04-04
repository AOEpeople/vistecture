package controller

import (
	"fmt"
	"os"
	"html/template"
	"../model/core"
	"../model/graphviz"
	"os/exec"
	"bytes"
	"strings"
)

type (
	DocumentationController struct {
		ProjectConfigPath *string
	}

	TemplateData struct {
		Project *core.Project
	}
)

func (DocumentationController DocumentationController) GraphvizAction(componentName string) {
	Project := loadProject(*DocumentationController.ProjectConfigPath)
	ProjectDrawer := graphviz.CreateProjectDrawer(Project)
	if componentName != ""  {
		Component, e := Project.FindComponent(componentName)
		if (e != nil) {
			fmt.Println(e)
			os.Exit(-1)
		}
		fmt.Print(ProjectDrawer.DrawComponent(&Component))
	} else {
		fmt.Print(ProjectDrawer.DrawComplete())
	}

}


func (DocumentationController DocumentationController) HTMLDocumentAction() {
	project := loadProject(*DocumentationController.ProjectConfigPath)
	tpl := template.New("htmldocument.tmpl")
	tpl.Funcs(template.FuncMap{
		"renderSVGInlineImage": func(Component core.Application) template.HTML {
			ProjectDrawer := graphviz.CreateProjectDrawer(project)
			stdInContent := ProjectDrawer.DrawComponent(&Component)

			commandName := "dot"
			dot := exec.Command(commandName,"-Tsvg")
			buf := new(bytes.Buffer)
			dot.Stdin = bytes.NewBufferString(stdInContent)
			dot.Stdout = buf
			e := dot.Run()
			if (e != nil) {
				fmt.Print(e)
			}
			dot.Wait()

			return template.HTML(buf.String())
		},
		"renderContent": func(content string) template.HTML {
			return template.HTML(strings.Replace(content," / ","<br />",-1))
		},
	})
	tpl, err := tpl.ParseFiles("templates/htmldocument.tmpl")
	if (err != nil) {
		fmt.Println(err)
		os.Exit(-1)
	}

	data := TemplateData {
		Project: project,
	}
	err = tpl.Execute(os.Stdout, data)
	if (err != nil) {
		fmt.Println(err)
		os.Exit(-1)
	}
}
