package main

import (
	"os"

	"github.com/AOEpeople/vistecture/controller"
	"gopkg.in/urfave/cli.v1"
)

func action(cb func()) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		cb()
		return nil
	}
}

func main() {
	var projectConfigPath, componentName, templatePath, iconPath, projectName, summaryRelation, hidePlanned string
	var skipValidation bool
	app := cli.NewApp()
	app.Name = "vistecture tool "
	app.Version = "0.5.3"
	app.Usage = "describing and analysing distributed or microservice-style architectures with its depenendcy!"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, definition",
			Value:       "project",
			Usage:       "Path to the project definition. Can be a file or a folder with json files",
			Destination: &projectConfigPath,
		},
		cli.StringFlag{
			Name:        "project, name, projectname",
			Value:       "None",
			Usage:       "Name of the project configuration",
			Destination: &projectName,
		},
		cli.BoolFlag{
			Name:        "skipValidation",
			Usage:       "Name of the project configuration",
			Destination: &skipValidation,
		},
	}

	analyzeController := controller.AnalyzeController{ProjectConfigPath: &projectConfigPath, ProjectName: &projectName}
	documentationController := controller.DocumentationController{ProjectConfigPath: &projectConfigPath, ProjectName: &projectName}

	app.Commands = []cli.Command{
		{
			Name:   "validate",
			Usage:  "Validates project JSON",
			Action: action(analyzeController.ValidateAction),
		},
		{
			Name:   "analyze",
			Usage:  "Analyses project structure. Detects cyclic dependencies etc",
			Action: action(analyzeController.AnalyzeAction),
		},
		{
			Name:   "documentation",
			Usage:  "Creates (living) documentation",
			Action: action(func() { documentationController.HTMLDocumentAction(templatePath, iconPath) }),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "templatePath",
					Value:       "templates/htmldocument.tmpl",
					Usage:       "Path of template that will be used",
					Destination: &templatePath,
				},
				cli.StringFlag{
					Name:        "iconPath",
					Value:       "templates/icons",
					Usage:       "Path of icons that will be in drawing components",
					Destination: &iconPath,
				},
			},
		},
		{
			Name:   "graph",
			Usage:  "Build graphviz format which can be used by dot or any other graphviz command. \n go run main.go graph | dot -Tpng -o graph.png \n See: http://www.graphviz.org/pdf/twopi.1.pdf",
			Action: action(func() { documentationController.GraphvizAction(componentName, iconPath, hidePlanned, skipValidation) }),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "application",
					Value:       "",
					Usage:       "Name of a application - then only a graph for this application will be drawn",
					Destination: &componentName,
				},
				cli.StringFlag{
					Name:        "iconPath",
					Value:       "templates/icons",
					Usage:       "Path of icons that will be in drawing components",
					Destination: &iconPath,
				},
				cli.StringFlag{
					Name:        "hidePlanned",
					Value:       "",
					Usage:       "Flag if planned applications should be drawn or not",
					Destination: &hidePlanned,
				},
			},
		},
		{
			Name:   "teamGraph",
			Usage:  "Build a overview of involved teams and the relations based from the architecture (Conways law)",
			Action: action(func() { documentationController.TeamGraphvizAction(summaryRelation) }),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "summaryRelation",
					Value:       "",
					Usage:       "if set then only one arrow is drawn between the teams",
					Destination: &summaryRelation,
				},
			},
		},
	}

	app.Run(os.Args)
}
