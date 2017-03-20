package graphviz

import (
	model "appdependency/model/core"
	"strings"
)

// EXTEND PROJECT

type ProjectDrawer struct {
	//inherit
	originalProject *model.Project
}

// Decorate Draw function
func (ProjectDrawer *ProjectDrawer) Draw() string {
	var result string
	result = "digraph { graph [] \n"
	// Nodes
	for key, componentList := range ProjectDrawer.originalProject.GetComponentsByGroup() {
		drawingResultInGroup := ""
		for _, component := range componentList {
			drawer := ComponentDrawer{originalComponent: component}
			drawingResultInGroup = drawingResultInGroup + drawer.Draw()
		}

		if key != model.NOGROUP {
			result = result + "subgraph \"cluster_"+key+"\" { label=\""+key+"\"; \n " + drawingResultInGroup + "\n}\n"
		} else {
			result = result  + drawingResultInGroup
		}
	}
	// Paths
	for _, component := range ProjectDrawer.originalProject.Components {


		// From components
		for _, dependency := range component.Dependencies {
			result = result + "\""+component.Name + "\" ->" + getGraphVizReference(dependency)+getEdgeLayoutFromDependency(dependency,component.Display)+"\n"
		}


		// From components/interfaces
		for _, providedInterface := range component.ProvidedServices {
			for _, dependency := range providedInterface.Dependencies {
				result = result + "\""+component.Name+"\":\""+providedInterface.Name+"\"->"+getGraphVizReference(dependency)+getEdgeLayoutFromDependency(dependency,component.Display)+"\n"
			}
		}
	}
	result = result + "}"
	return result
}


//Get GRaphviz style reference
func getGraphVizReference(Dependency model.Dependency) string {
	if (strings.Contains(Dependency.Reference,".")) {
		splitted := strings.Split(Dependency.Reference, ".")
		return "\""+splitted[0]+"\":\""+splitted[1]+"\"";
	}
	return "\""+Dependency.Reference+"\"";
}

func getEdgeLayoutFromDependency(dependency model.Dependency, display model.ComponentDisplaySettings) string {

	edgeLayout := "["
	if (display.BorderColor != "") {
		edgeLayout =edgeLayout+"color=\""+display.BorderColor+"\""
	} else {
		edgeLayout =edgeLayout+"color=\"#333333\""
	}
	edgeLayout =edgeLayout+", fontsize=\"10\", fontcolor=\"#555555\" "
	if dependency.Relationship == "acl" {
		edgeLayout = edgeLayout + ", dir=both, arrowtail=\"box\", taillabel=<<font color=\"red\" ><b>acl</b></font>>"
	} else {
		edgeLayout = edgeLayout + ", label=\""+dependency.Relationship+"\""
	}
	if dependency.Relationship == "customer-supplier" || dependency.Relationship == "conformist" {
		edgeLayout = edgeLayout + "weight=3, style=\"bold\""
	} else if dependency.IsBrowserBased {
		edgeLayout = edgeLayout + "style=\"dashed\""
	}
	if dependency.IsSameLevel {
		edgeLayout = edgeLayout + ", constraint=false"
	}
	return edgeLayout + "]"
}


// Factory
func CreateProjectDrawer(Project *model.Project) *ProjectDrawer {
	var Drawer ProjectDrawer
	Drawer.originalProject = Project
	return &Drawer
}



// EXTEND COMPONENT
type ComponentDrawer struct {
	//inherit
	originalComponent *model.Component
}

// Decorate Draw function
func (ComponentDrawer ComponentDrawer) Draw() string {
	var result string
	Component := ComponentDrawer.originalComponent

	icon := ""
	switch Component.Technology {
	case "go","scala","php", "anypoint", "akeneo", "magento","keycloak":
		icon = "<IMG SRC=\"templates/res/"+Component.Technology+".png\" scale=\"true\"/>"
	}

	tableHeaderColor := ""
	switch Component.Category {
	case "external":
		tableHeaderColor = "#8e0909"
	default:
		tableHeaderColor = "#1B4E5E"

	}
	// see http://www.graphviz.org/doc/info/shapes.html
	// see http://4webmaster.de/wiki/Graphviz-Tutorial#Die_Darstellung_von_Edges_ver.C3.A4ndern
	result = result + "\""+Component.Name+"\" [shape=plaintext "
	if (Component.Display.BorderColor != "") {
		result = result + ", color=\""+Component.Display.BorderColor+"\""
	}



	result = result + ", label=<<TABLE BGCOLOR=\"#1B4E5E\" ROWS=\"*\" CELLPADDING=\"3\" BORDER=\"2\" CELLBORDER=\"0\" CELLSPACING=\"0\"> \n"
	result = result+" <TR ><TD BGCOLOR=\""+tableHeaderColor+"\"><FONT COLOR=\"#fefefe\">"+strings.Replace(strings.ToTitle(Component.Name),"/","\n<BR />",1)+"</FONT></TD><TD BGCOLOR=\""+tableHeaderColor+"\" width=\"50\" height=\"30\" fixedsize=\"true\" >"+icon+"</TD></TR> \n"
	if Component.Description != "" {
		result = result+" <TR ><TD COLSPAN=\"2\" BGCOLOR=\"#aaaaaa\"><FONT POINT-SIZE=\"10\">"+Component.Description+"</FONT></TD></TR> \n"
	}
	for _, service := range Component.ProvidedServices {
		var color string
		switch service.Type {
		case "api":
			color = "#A3C7D4"
		case "gui":
			color = "#D4C1E0"
		case "exchange":
			color = "#BEE8D2"
		default:
			color = "#CFCFCF"
		}
		result = result + "<TR><TD COLSPAN=\"2\"  align=\"CENTER\" PORT=\""+service.Name+"\" BGCOLOR=\""+color+"\"><FONT POINT-SIZE=\"10\">"+service.Type+":"+service.Name+"</FONT></TD></TR>"
	}
	result = result + "</TABLE>>];\n"
	return result

}