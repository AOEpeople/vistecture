package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/AOEpeople/vistecture/v2/application"
	"github.com/AOEpeople/vistecture/v2/controller"
	"github.com/AOEpeople/vistecture/v2/controller/web"
	"github.com/AOEpeople/vistecture/v2/model/core"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
)

type (
	projectInjectAble interface {
		Inject(*core.Project)
	}
)

var (
	//global cli flags
	projectConfigFile, projectSubViewName string
	skipValidation                        bool
	//server cli flags
	serverPort            int
	localTemplateFolder   string
	staticDocumentsFolder string
)

func actionFunc(lazyProjectInjectAble projectInjectAble, cb func()) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		project := loadProject(projectConfigFile, projectSubViewName, skipValidation)
		lazyProjectInjectAble.Inject(project)
		cb()
		return nil
	}
}

func main() {
	var componentName, templatePath, iconPath, summaryRelation, hidePlanned string

	app := cli.NewApp()
	app.Name = "vistecture tool "
	app.Version = "2.0.11"
	app.Usage = "describing and analysing distributed or microservice-style architectures with its depenendcies."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config",
			Value:       "",
			Usage:       "Path to the project config file",
			Destination: &projectConfigFile,
		},
		cli.StringFlag{
			Name:        "subview",
			Value:       "",
			Usage:       "Name of the projects subview - if you want to limit the action to subview",
			Destination: &projectSubViewName,
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
			Action: validate,
		},
		{
			Name:   "list",
			Usage:  "lists the apps",
			Action: listApps,
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
			Name:   "groupGraph",
			Usage:  "Build graphviz format that shows only the group of services and its dependencies.",
			Action: actionFunc(documentationController, func() { documentationController.GroupGraphvizAction(summaryRelation) }),
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
				cli.StringFlag{
					Name:        "localTemplateFolder",
					Value:       "",
					Usage:       "if set then this template folder will be used to serve the view - otherwise a standard template gets loaded automatically",
					Destination: &localTemplateFolder,
				},
				cli.StringFlag{
					Name:        "staticDocumentsFolder",
					Value:       "",
					Usage:       "if set then this  folder will be scanned for files that are linked in the mainmenu then",
					Destination: &staticDocumentsFolder,
				},
			},
		},
	}

	_ = app.Run(os.Args)
}

func loadProject(configFile string, subViewName string, skipValidation bool) *core.Project {
	loader := application.ProjectLoader{StrictMode: !skipValidation}
	project, err := loader.LoadProjectFromConfigFile(configFile, subViewName)

	if err != nil {
		log.Println(err)
		if !skipValidation {
			log.Fatal("project loading aborted.")
		}
	}
	return project
}

func validate(_ *cli.Context) error {
	loader := application.ProjectLoader{StrictMode: !skipValidation}
	project, err := loader.LoadProjectFromConfigFile(projectConfigFile, projectSubViewName)

	if err != nil {
		log.Println(err)
	}
	var validationErrors []error
	if project != nil {
		validationErrors = project.Validate()
		for _, valErr := range validationErrors {
			log.Println(valErr)
		}
	}
	if err != nil || len(validationErrors) > 0 {
		log.Fatal("Not valid")
	} else {
		log.Println("valid")
	}
	return nil
}

func listApps(_ *cli.Context) error {
	project := loadProject(projectConfigFile, projectSubViewName, true)
	for _, app := range project.Applications {
		log.Printf("Name: %v Id: %v", app.Name, app.Id)
	}
	return nil
}

func startServer(_ *cli.Context) error {
	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%v", serverPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	webProjectController := web.ProjectController{}
	loader := application.ProjectLoader{}
	definitions, err := loader.LoadProjectConfig(projectConfigFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	webProjectController.Inject(definitions, &loader, path.Dir(projectConfigFile), skipValidation)

	// This will serve files under http://localhost:8000/documents/<filename>
	r.PathPrefix("/documents/").Handler(http.StripPrefix("/documents/", http.FileServer(http.Dir(staticDocumentsFolder))))

	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		webProjectController.DataAction(w, r, staticDocumentsFolder)
	})

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webProjectController.IndexAction(w, r, localTemplateFolder)
	})

	log.Printf("Starting server:%v \n", serverPort)
	log.Fatal(srv.ListenAndServe())
	return nil
}
