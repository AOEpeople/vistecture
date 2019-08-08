package graphviz

import (
	"fmt"
	"strings"

	model "github.com/AOEpeople/vistecture/v2/model/core"
)

// EXTEND PROJECT

type (
	GroupDrawer struct {
		project             *model.Project
		summaryRelationOnly bool
	}
	groupRelation struct {
		Relationship   string
		ToGroup         string
		ForApplication string
		FromApplication string
	}
)

// Factory
func CreateGroupDrawer(Project *model.Project, summaryRelation bool) *GroupDrawer {
	var Drawer GroupDrawer
	Drawer.project = Project
	Drawer.summaryRelationOnly = summaryRelation
	return &Drawer
}

func (d *GroupDrawer) DrawComplete() string {

	groups := make(map[string][]*model.Application)
	groupOutgoing := make(map[string][]groupRelation)

	// Build Graph Infos
	for _, application := range d.project.Applications {
		groupName := application.GetMainGroup()
		if groupName == "" {
			groupName = "UNGROUPED"
		}

		groups[groupName] = append(groups[groupName], application)
		dependencies := application.GetAllDependencies()
		for _, dependency := range dependencies {
			dependencyApplication, e := dependency.GetApplication(d.project)
			if e != nil {
				continue
			}
			depGroupName := dependencyApplication.GetMainGroup()
			if depGroupName == "" {
				depGroupName = "UNGROUPED"
			}
			if depGroupName == groupName {
				continue
			}

			relationShip := dependency.Relationship
			if relationShip == "" {
				if dependencyApplication.IsOpenHostApp() {
					relationShip = "open-host"
				}
			}
			referenceToApplicationAlreadyPresent := false
			for _, rel := range groupOutgoing[application.Team] {
				if rel.ForApplication == dependencyApplication.Name && rel.FromApplication == application.Name {
					referenceToApplicationAlreadyPresent = true
				}
			}
			if referenceToApplicationAlreadyPresent {
				continue
			}
			groupOutgoing[groupName] = append(groupOutgoing[groupName],
				groupRelation{
					Relationship:   relationShip,
					ToGroup:        depGroupName,
					ForApplication: dependencyApplication.Name,
					FromApplication: application.Name,
				})

		}
	}

	//Draw Graph
	colors := []string{"#9013a0", "#2936c4", "#147724", "#22a398", "#9e8142", "#bcae67", "#d62a2a"}
	color := "#333333"
	var result string
	result = "digraph { graph [overlap=false] \n"
	i := 0
	for group, applications := range groups {
		i++
		if (len(colors)) <= i {
			i = 0
		}
		color = colors[i]

		result = result + d.DrawGroup(group, applications, color) + "\n"
		if d.summaryRelationOnly {
			//Draw relation to team only
			strongestToGroup := make(map[string]string)
			for _, relation := range groupOutgoing[group] {
				if currentRelation, ok := strongestToGroup[relation.ToGroup]; ok {
					if isStrongerRelation(currentRelation, relation.Relationship) {
						strongestToGroup[relation.ToGroup] = relation.Relationship
					}
				} else {
					strongestToGroup[relation.ToGroup] = relation.Relationship
				}
			}
			//Draw relation to every application
			for toGroup, relationshipType := range strongestToGroup {
				edgeLayout := edgeLayout(relationshipType)
				edgeLayout += ", label=\"" + relationshipType + "\""
				result = result + "\"" + group + "\"->\"" + toGroup + "\"[color=\"" + color + "\" " + edgeLayout + "]\n"
			}

		} else {
			//Draw relation to every application
			for _, relation := range groupOutgoing[group] {
				edgeLayout := edgeLayout(relation.Relationship)

				result = result + "\"" + group + "\":\"" + relation.FromApplication + "\"->\"" + relation.ToGroup + "\":\"" + relation.ForApplication + "\"[color=\"" + color + "\" " + edgeLayout + "]\n"
			}
		}

	}

	result = result + "}"
	return result
}

func (d *GroupDrawer) DrawGroup(group string, applications []*model.Application, tableHeaderColor string) string {
	var result string

	// see http://www.graphviz.org/doc/info/shapes.html
	// see http://4webmaster.de/wiki/Graphviz-Tutorial#Die_Darstellung_von_Edges_ver.C3.A4ndern
	result += "\"" + group + "\" [shape=plaintext "

	result += ", label=<<TABLE BGCOLOR=\"#1B4E5E\" ROWS=\"*\" CELLPADDING=\"3\" BORDER=\"2\" CELLBORDER=\"0\" CELLSPACING=\"0\"> \n"
	result += " <TR ><TD BGCOLOR=\"" + tableHeaderColor + "\"><FONT COLOR=\"#fefefe\">" + strings.Replace(strings.ToTitle(group), " / ", "\n<BR />", 1) + "</FONT></TD></TR> \n"

	for _, app := range applications {
		color := "#CFCFCF"

		teamName := ""
		if app.Team != "" {
			teamName = fmt.Sprintf(" (%v)",escape(app.Team))
		}
		result += "<TR><TD COLSPAN=\"2\"  align=\"CENTER\" PORT=\"" + escape(app.Name) + "\" BGCOLOR=\"" + color + "\">"
		result += "<FONT POINT-SIZE=\"10\">" + escape(app.Name) + teamName + "</FONT>"

		result += "</TD></TR>"
	}
	result += "</TABLE>>];\n"
	return result
}
