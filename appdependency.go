package main

import (
	"os"
	"gopkg.in/urfave/cli.v1"
	controller "appdependency/controller"
)


func main() {

	var projectConfigPath string

	app := cli.NewApp()
	app.Name = "appdependency tool "
	app.Usage = " describing and analysing distributed or microservice-style architectures with its depenendcy!"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "config",
			Value: "project",
			Usage: "Path to the project definition. Can be a file or a folder with json files",
			Destination: &projectConfigPath,
		},
	}


	app.Commands = []cli.Command{
		{
			Name:    "validate",
			Usage:   "Validates project JSON",
			Action:  func(c *cli.Context) error {
				AnalyzeController := controller.AnalyzeController{projectConfigPath}
				AnalyzeController.ValidateAction()
				return nil
			},
		},
		{
			Name:    "analyze",
			Usage:   "Analyses project structure. Detects cyclic dependencies etc",
			Action:  func(c *cli.Context) error {
				AnalyzeController := controller.AnalyzeController{projectConfigPath}
				AnalyzeController.AnalyzeAction()
				return nil
			},
		},
		{
			Name:    "documentation",
			Usage:   "Creates (living) documentation",
			Action:  func(c *cli.Context) error {
				DocumentationController := controller.DocumentationController{projectConfigPath}
				DocumentationController.HTMLDocumentAction()

				return nil
			},
		},
		{
			Name:    "graph",
			Usage:   "Build graphviz format which can be used by dot or any other graphviz command. \n go run main.go graph | dot -Tpng -o graph.png \n See: http://www.graphviz.org/pdf/twopi.1.pdf",
			Action:  func(c *cli.Context) error {
				DocumentationController := controller.DocumentationController{projectConfigPath}
				DocumentationController.GraphvizAction()
				return nil
			},
		},
	}

	app.Run(os.Args)


}