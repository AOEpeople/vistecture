package main

import (
	"appdependency/controller"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func action(cb func()) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		cb()
		return nil
	}
}

func main() {
	var projectConfigPath, componentName string

	app := cli.NewApp()
	app.Name = "appdependency tool "
	app.Usage = " describing and analysing distributed or microservice-style architectures with its depenendcy!"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config",
			Value:       "project",
			Usage:       "Path to the project definition. Can be a file or a folder with json files",
			Destination: &projectConfigPath,
		},
	}

	analyzeController := controller.AnalyzeController{ProjectConfigPath: &projectConfigPath}
	documentationController := controller.DocumentationController{ProjectConfigPath: &projectConfigPath}

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
			Action: action(documentationController.HTMLDocumentAction),
		},
		{
			Name:   "graph",
			Usage:  "Build graphviz format which can be used by dot or any other graphviz command. \n go run main.go graph | dot -Tpng -o graph.png \n See: http://www.graphviz.org/pdf/twopi.1.pdf",
			Action: action(func() { documentationController.GraphvizAction(componentName) }),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "component",
					Value:       "",
					Usage:       "Name of a component - then only a graph for this component will be drawn",
					Destination: &componentName,
				},
			},
		},
	}

	app.Run(os.Args)
}
