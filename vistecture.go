package main

import (
	"log"
	"net/http"
	"os"

	"fmt"

	"github.com/AOEpeople/vistecture/application"
	"github.com/AOEpeople/vistecture/controller"
	"github.com/AOEpeople/vistecture/controller/web"
	"github.com/AOEpeople/vistecture/model/core"
	"gopkg.in/urfave/cli.v1"
)

type (
	projectInjectAble interface {
		Inject(*core.Project)
	}
)

var (
	//global cli flags
	projectConfigPath, projectName string
	skipValidation                 bool
	//server cli flag
	serverPort int
)

func actionFunc(lazyProjectInjectAble projectInjectAble, cb func()) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		project := loadProject(projectConfigPath, projectName, skipValidation)
		lazyProjectInjectAble.Inject(project)
		cb()
		return nil
	}
}

func main() {
	var componentName, templatePath, iconPath, summaryRelation, hidePlanned string

	app := cli.NewApp()
	app.Name = "vistecture tool "
	app.Version = "1.0.0"
	app.Usage = "describing and analysing distributed or microservice-style architectures with its depenendcies."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, definition",
			Value:       "project",
			Usage:       "Path to the project definition. Can be a file or a folder with json files",
			Destination: &projectConfigPath,
		},
		cli.StringFlag{
			Name:        "project, name, projectname",
			Value:       "",
			Usage:       "Name of the project configuration to use. If not set all definied applications will be used",
			Destination: &projectName,
		},
		cli.BoolFlag{
			Name:        "skipValidation",
			Usage:       "Skip the validation of the project",
			Destination: &skipValidation,
		},
	}

	analyzeController := &controller.AnalyzeController{}
	documentationController := &controller.DocumentationController{}

	app.Commands = []cli.Command{
		{
			Name:   "validate",
			Usage:  "Validates project JSON",
			Action: actionFunc(analyzeController, func() { log.Println("valid") }),
		},
		{
			Name:   "analyze",
			Usage:  "Analyses project structure. Detects cyclic dependencies etc",
			Action: actionFunc(analyzeController, analyzeController.AnalyzeAction),
		},
		{
			Name:   "documentation",
			Usage:  "Creates (living) documentation",
			Action: actionFunc(documentationController, func() { documentationController.HTMLDocumentAction(templatePath, iconPath) }),
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
			Action: actionFunc(documentationController, func() { documentationController.GraphvizAction(componentName, iconPath, hidePlanned, skipValidation) }),
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
			Action: actionFunc(documentationController, func() { documentationController.TeamGraphvizAction(summaryRelation) }),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "summaryRelation",
					Value:       "",
					Usage:       "if set then only one arrow is drawn between the teams",
					Destination: &summaryRelation,
				},
			},
		},
		{
			Name:   "serve",
			Usage:  "Runs the vistecture webserver",
			Action: startServer,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "port",
					Value:       8080,
					Usage:       "set the port (default 8080)",
					Destination: &serverPort,
				},
			},
		},
	}

	app.Run(os.Args)
}

func loadProject(ProjectConfigPath string, ProjectName string, skipValidation bool) *core.Project {
	project, errors := application.CreateProjectByName(ProjectConfigPath, ProjectName, skipValidation)

	if len(errors) > 0 {
		for _, err := range errors {
			log.Print("project creation failed because of: ", err)
		}
		log.Fatal("project loading aborted.")
	}
	return project
}

func startServer(c *cli.Context) error {
	webProjectController := web.ProjectController{}
	definitions, err := application.LoadProjectDefinitions(projectConfigPath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	webProjectController.Inject(definitions)

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {

		webProjectController.DataAction(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		webProjectController.IndexAction(w, r)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", serverPort), nil))
	return nil
}
