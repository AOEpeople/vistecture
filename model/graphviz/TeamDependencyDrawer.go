package graphviz

import (
	"strings"

	model "github.com/AOEpeople/vistecture/v2/model/core"
)

// EXTEND PROJECT

type (
	TeamDependencyDrawer struct {
		project             *model.Project
		summaryRelationOnly bool
	}
	OutgoingTeamRelation struct {
		Relationship   string
		ToTeam         string
		ForApplication string
	}
)

// Factory
func CreateTeamDependencyDrawer(Project *model.Project, summaryRelationOnly bool) *TeamDependencyDrawer {
	var Drawer TeamDependencyDrawer
	Drawer.project = Project
	Drawer.summaryRelationOnly = summaryRelationOnly
	return &Drawer
}

// Decorate Draw function
func (d *TeamDependencyDrawer) DrawComplete() string {

	teams := make(map[string][]*model.Application)
	teamOutgoing := make(map[string][]OutgoingTeamRelation)

	// Build Graph Infos
	for _, application := range d.project.Applications {
		if application.Team == "" {
			continue
		}
		teams[application.Team] = append(teams[application.Team], application)
		dependencies := application.GetAllDependencies()
		for _, dependency := range dependencies {
			dependencyApplication, e := dependency.GetApplication(d.project)
			if e != nil {
				continue
			}
			if dependencyApplication.Team == application.Team {
				continue
			}
			if dependencyApplication.Team != "" {
				relationShip := dependency.Relationship
				if relationShip == "" {
					if dependencyApplication.IsOpenHostApp() {
						relationShip = "open-host"
					}
				}
				referenceToApplicationAlreadyPResent := false
				for _, rel := range teamOutgoing[application.Team] {
					if rel.ForApplication == dependencyApplication.Name {
						referenceToApplicationAlreadyPResent = true
					}
				}
				if referenceToApplicationAlreadyPResent {
					continue
				}
				teamOutgoing[application.Team] = append(teamOutgoing[application.Team],
					OutgoingTeamRelation{
						Relationship:   relationShip,
						ToTeam:         dependencyApplication.Team,
						ForApplication: dependencyApplication.Name,
					})
			}
		}
	}

	//Draw Graph
	colors := []string{"#9013a0", "#2936c4", "#147724", "#22a398", "#9e8142", "#bcae67", "#d62a2a"}
	color := "#333333"
	var result string
	result = "digraph { graph [overlap=false] \n"
	i := 0
	for team, applications := range teams {
		i++
		if (len(colors)) <= i {
			i = 0
		}
		color = colors[i]

		result = result + d.DrawTeam(team, applications, color) + "\n"
		if d.summaryRelationOnly {
			//Draw relation to team only
			strongestToTeam := make(map[string]string)
			for _, relation := range teamOutgoing[team] {
				if currentRelation, ok := strongestToTeam[relation.ToTeam]; ok {
					if isStrongerRelation(currentRelation, relation.Relationship) {
						strongestToTeam[relation.ToTeam] = relation.Relationship
					}
				} else {
					strongestToTeam[relation.ToTeam] = relation.Relationship
				}
			}
			//Draw relation to every application
			for toTeam, relationshipType := range strongestToTeam {
				edgeLayout := edgeLayout(relationshipType)
				edgeLayout += ", label=\"" + relationshipType + "\""
				result = result + "\"" + team + "\"->\"" + toTeam + "\"[color=\"" + color + "\" " + edgeLayout + "]\n"
			}

		} else {
			//Draw relation to every application
			for _, relation := range teamOutgoing[team] {
				edgeLayout := edgeLayout(relation.Relationship)

				result = result + "\"" + team + "\"->\"" + relation.ToTeam + "\":\"" + relation.ForApplication + "\"[color=\"" + color + "\" " + edgeLayout + "]\n"
			}
		}

	}

	result = result + "}"
	return result
}
func edgeLayout(relationShipType string) string {
	edgeLayout := ""
	if relationShipType == "acl" {
		edgeLayout += ", style=\"dashed\""
	}
	if relationShipType == "open-host" {
		edgeLayout += ", style=\"dashed\""
	}
	if relationShipType == "customer-supplier" {
		edgeLayout += ", weight=2, style=\"bold\""
	}
	if relationShipType == "conformist" || relationShipType == "partnership" {
		edgeLayout += ", weight=3, style=\"bold\""
	}

	return edgeLayout
}

// Decorate Draw function - Draws only a component with its direct dependencies and direct callers
func (d *TeamDependencyDrawer) DrawTeam(team string, applications []*model.Application, tableHeaderColor string) string {
	var result string

	if tableHeaderColor == "" {
		tableHeaderColor = "#333333"
	}

	// see http://www.graphviz.org/doc/info/shapes.html
	// see http://4webmaster.de/wiki/Graphviz-Tutorial#Die_Darstellung_von_Edges_ver.C3.A4ndern
	result += "\"" + team + "\" [shape=plaintext "

	result += ", label=<<TABLE BGCOLOR=\"#1B4E5E\" ROWS=\"*\" CELLPADDING=\"3\" BORDER=\"2\" CELLBORDER=\"0\" CELLSPACING=\"0\"> \n"
	result += " <TR ><TD BGCOLOR=\"" + tableHeaderColor + "\"><FONT COLOR=\"#fefefe\">" + strings.Replace(strings.ToTitle(team), " / ", "\n<BR />", 1) + "</FONT></TD></TR> \n"

	for _, app := range applications {
		color := "#CFCFCF"

		result += "<TR><TD COLSPAN=\"2\"  align=\"CENTER\" PORT=\"" + escape(app.Name) + "\" BGCOLOR=\"" + color + "\">"
		result += "<FONT POINT-SIZE=\"10\">" + escape(app.Name) + "</FONT>"

		result += "</TD></TR>"
	}
	result += "</TABLE>>];\n"
	return result
}

func isStrongerRelation(current string, toCheck string) bool {
	if relationTypeWeight(current) < relationTypeWeight(toCheck) {
		return true
	}
	return false
}

func relationTypeWeight(current string) int {
	if current == "acl" {
		return 0
	}
	if current == "customer-supplier" {
		return 2
	}

	if current == "conformist" {
		return 3
	}

	if current == "partnership" {
		return 4
	}
	return 1
}
