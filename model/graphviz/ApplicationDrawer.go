package graphviz

import (
	"os"
	"strings"

	model "github.com/AOEpeople/vistecture/model/core"
)

// EXTEND COMPONENT
type ApplicationDrawer struct {
	//inherit
	originalComponent *model.Application

	iconPath string
}

// Decorate Draw function
func (ComponentDrawer ApplicationDrawer) Draw() string {
	var result string
	Component := ComponentDrawer.originalComponent
	icon := ""
	iconPath := ComponentDrawer.iconPath + "/" + strings.ToLower(Component.Technology) + ".png"
	if _, err := os.Stat(iconPath); err == nil {
		icon = "<IMG SRC=\"" + iconPath + "\" scale=\"true\"/>"
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

	if Component.Status == model.STATUS_PLANNED {
		result += "\"" + Component.Name + "\" [xlabel=\"planned\", shape=plaintext "
		if Component.Display.BorderColor != "" {
			result += ", color=\"#BBBBBB\""
		}
		result += ", label=<<TABLE BGCOLOR=\"#1B4E5E\" ROWS=\"*\" CELLPADDING=\"3\" BORDER=\"0\" CELLBORDER=\"0\" CELLSPACING=\"0\"> \n"
		tableHeaderColor = "#BBBBBB"
	} else {
		result += "\"" + Component.Name + "\" [style=dotted, shape=plaintext "
		if Component.Display.BorderColor != "" {
			result += ", color=\"" + Component.Display.BorderColor + "\""
		}
		result += ", label=<<TABLE BGCOLOR=\"#1B4E5E\" ROWS=\"*\" CELLPADDING=\"3\" BORDER=\"2\" CELLBORDER=\"0\" CELLSPACING=\"0\"> \n"
	}

	result += " <TR ><TD BGCOLOR=\"" + tableHeaderColor + "\"><FONT COLOR=\"#fefefe\">" + strings.Replace(strings.ToTitle(Component.Name), " / ", "\n<BR />", 1) + "</FONT></TD><TD BGCOLOR=\"" + tableHeaderColor + "\" width=\"50\" height=\"30\" fixedsize=\"true\" >" + icon + "</TD></TR> \n"
	if Component.Title != "" {
		result += " <TR ><TD COLSPAN=\"2\" BGCOLOR=\"#aaaaaa\"><FONT POINT-SIZE=\"10\">" + escape(Component.Title) + "</FONT></TD></TR> \n"
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
		result += "<TR><TD COLSPAN=\"2\"  align=\"CENTER\" PORT=\"" + escape(service.Name) + "\" BGCOLOR=\"" + color + "\">"
		result += "<FONT POINT-SIZE=\"10\">" + service.Type + ":" + escape(service.Name) + "</FONT>"
		if service.IsOpenHost {
			result += " <FONT COLOR=\"#33911a\">♡</FONT>"
		}
		result += "</TD></TR>"
	}
	result += "</TABLE>>];\n"
	return result

}

func escape(value string) string {

	value = strings.Replace(value, "&", "&amp;", -1)

	return value
}
