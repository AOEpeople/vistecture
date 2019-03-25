package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"log"

	"github.com/AOEpeople/vistecture/v2/model/core"
	"github.com/AOEpeople/vistecture/v2/model/graphviz"
)

type (
	DocumentationController struct {
		project *core.Project
	}

	TemplateData struct {
		Project *core.Project
	}
)

func (d *DocumentationController) Inject(project *core.Project) {
	d.project = project
}

func (d *DocumentationController) GraphvizAction(componentName string, iconPath string, hidePlanned string, skipValidation bool) {
	projectDrawer := graphviz.CreateProjectDrawer(d.project, iconPath)
	if componentName != "" {
		Component, e := d.project.FindApplication(componentName)
		if e != nil {
			log.Println(e)
			if !skipValidation {
				os.Exit(-1)
			}
		}
		if Component == nil {
			os.Exit(-1)
		}
		fmt.Print(projectDrawer.DrawComponent(Component))
	} else {
		fmt.Print(projectDrawer.DrawComplete(hidePlanned == "1"))
	}

}

func (d *DocumentationController) TeamGraphvizAction(summaryRelation string) {
	drawer := graphviz.CreateTeamDependencyDrawer(d.project, summaryRelation != "")
	fmt.Print(drawer.DrawComplete())
}

func (d *DocumentationController) HTMLDocumentAction(templatePath string, iconPath string) {
	tpl := template.New(filepath.Base(templatePath))

	tpl.Funcs(template.FuncMap{
		"renderSVGInlineImage": func(Component core.Application) template.HTML {
			ProjectDrawer := graphviz.CreateProjectDrawer(d.project, iconPath)
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
		Project: d.project,
	}
	err = tpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
