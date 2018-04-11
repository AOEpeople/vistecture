package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"path/filepath"

	"github.com/AOEpeople/vistecture/model/core"
	"github.com/AOEpeople/vistecture/model/graphviz"
)

type (
	DocumentationController struct {
		ProjectConfigPath *string
		ProjectName *string
	}

	TemplateData struct {
		Project *core.Project
	}
)

func (DocumentationController DocumentationController) GraphvizAction(componentName string, iconPath string) {
	Project := loadProject(*DocumentationController.ProjectConfigPath, *DocumentationController.ProjectName)
	ProjectDrawer := graphviz.CreateProjectDrawer(Project, iconPath)
	if componentName != "" {
		Component, e := Project.FindApplication(componentName)
		if e != nil {
			fmt.Println(e)
			os.Exit(-1)
		}
		fmt.Print(ProjectDrawer.DrawComponent(&Component))
	} else {
		fmt.Print(ProjectDrawer.DrawComplete())
	}

}

func (DocumentationController DocumentationController) HTMLDocumentAction(templatePath string, iconPath string) {
	project := loadProject(*DocumentationController.ProjectConfigPath, *DocumentationController.ProjectName)
	tpl := template.New(filepath.Base(templatePath))

	tpl.Funcs(template.FuncMap{
		"renderSVGInlineImage": func(Component core.Application) template.HTML {
			ProjectDrawer := graphviz.CreateProjectDrawer(project, iconPath)
			stdInContent := ProjectDrawer.DrawComponent(&Component)

			commandName := "dot"
			dot := exec.Command(commandName, "-Tsvg")
			buf := new(bytes.Buffer)
			dot.Stdin = bytes.NewBufferString(stdInContent)
			dot.Stdout = buf
			e := dot.Run()
			if e != nil {
				fmt.Print(e)
			}
			dot.Wait()

			return template.HTML(buf.String())
		},
		"renderContent": func(content string) template.HTML {
			return template.HTML(strings.Replace(content, " / ", "<br />", -1))
		},
	})
	tpl, err := tpl.ParseFiles(templatePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	data := TemplateData{
		Project: project,
	}
	err = tpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
