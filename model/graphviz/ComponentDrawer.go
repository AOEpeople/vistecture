package graphviz

import (
	"strings"
	model "appdependency/model/core"
)


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
	case "go", "scala", "php", "anypoint", "akeneo", "magento", "keycloak":
		icon = "<IMG SRC=\"templates/res/" + Component.Technology + ".png\" scale=\"true\"/>"
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
	result += "\"" + Component.Name + "\" [shape=plaintext "
	if Component.Display.BorderColor != "" {
		result += ", color=\"" + Component.Display.BorderColor + "\""
	}

	result += ", label=<<TABLE BGCOLOR=\"#1B4E5E\" ROWS=\"*\" CELLPADDING=\"3\" BORDER=\"2\" CELLBORDER=\"0\" CELLSPACING=\"0\"> \n"
	result += " <TR ><TD BGCOLOR=\"" + tableHeaderColor + "\"><FONT COLOR=\"#fefefe\">" + strings.Replace(strings.ToTitle(Component.Name), " / ", "\n<BR />", 1) + "</FONT></TD><TD BGCOLOR=\"" + tableHeaderColor + "\" width=\"50\" height=\"30\" fixedsize=\"true\" >" + icon + "</TD></TR> \n"
	if Component.Description != "" {
		result += " <TR ><TD COLSPAN=\"2\" BGCOLOR=\"#aaaaaa\"><FONT POINT-SIZE=\"10\">" + renderDescription(Component.GetSummary()) + "</FONT></TD></TR> \n"
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
		result += "<TR><TD COLSPAN=\"2\"  align=\"CENTER\" PORT=\"" + service.Name + "\" BGCOLOR=\"" + color + "\">"
		result += "<FONT POINT-SIZE=\"10\">"+service.Type + ":" + service.Name + "</FONT>"
		if service.Description != "" {
			result += "<BR /><FONT POINT-SIZE=\"8\">"+service.Description + "</FONT>"
		}

		result += "</TD></TR>"
	}
	result += "</TABLE>>];\n"
	return result

}

func renderDescription(description string) string {
	description = strings.Replace(description, " / ","<BR /> ",-1)
	return description
}
