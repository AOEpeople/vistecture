package graphviz

import (
	"strings"

	model "github.com/AOEpeople/vistecture/model/core"
)

// EXTEND PROJECT

type ProjectDrawer struct {
	//inherit
	originalProject *model.Project
	iconPath        string
}

// Decorate Draw function
func (projectDrawer *ProjectDrawer) DrawComplete(hidePlanned bool) string {
	var result string
	result = "digraph { graph [bgcolor=\"transparent\",overlap=false] \n"
	// Nodes
	result += projectDrawer.drawGroups(projectDrawer.originalProject.GetApplicationsRootGroup(), hidePlanned)

	// Paths
	for _, component := range projectDrawer.originalProject.Applications {
		if component.Status == model.STATUS_PLANNED && hidePlanned {
			continue
		}
		result = result + projectDrawer.drawComponentOutgoingRelations(component, hidePlanned)
	}
	result = result + "}"
	return result
}

//drawGroups - draws Nodes in the group and recursive calls drawGroup for subGroups
func (projectDrawer *ProjectDrawer) drawGroups(appsByGroup *model.ApplicationsByGroup, hidePlanned bool) string {
	result := ""

	for _, subGroup := range appsByGroup.SubGroups {
		result += projectDrawer.drawGroups(subGroup, hidePlanned)
	}

	for _, component := range appsByGroup.Applications {
		if component.Status == model.STATUS_PLANNED && hidePlanned {
			continue
		}
		drawer := ApplicationDrawer{originalComponent: component, iconPath: projectDrawer.iconPath}
		result += drawer.Draw(hidePlanned)
	}
	if !appsByGroup.IsRoot {
		result = "subgraph \"cluster_" + appsByGroup.GroupName + "\" { label=\"" + appsByGroup.GroupName + "\"; \n " + result + "\n}\n"
	}

	return result
}

// Decorate Draw function - Draws only a component with its direct dependencies and direct callers
func (ProjectDrawer *ProjectDrawer) DrawComponent(Component *model.Application) string {
	var result string
	result = "digraph { graph [] \n"
	drawer := ApplicationDrawer{originalComponent: Component, iconPath: ProjectDrawer.iconPath}
	result = result + drawer.Draw(false)

	// Draw outgoing:
	result = result + ProjectDrawer.drawComponentOutgoingRelations(Component, false)
	allRelatedComponents, _ := Component.GetAllDependencyApplications(ProjectDrawer.originalProject)
	for _, relatedComponent := range allRelatedComponents {
		drawer := ApplicationDrawer{originalComponent: relatedComponent, iconPath: ProjectDrawer.iconPath}
		result = result + drawer.Draw(false)
	}
	//Draw incoming

	allDependendComponents := ProjectDrawer.originalProject.FindApplicationThatReferenceTo(Component, false)
	for _, relatedComponent := range allDependendComponents {
		drawer := ApplicationDrawer{originalComponent: relatedComponent, iconPath: ProjectDrawer.iconPath}
		result = result + drawer.Draw(false)
		dependencies, e := relatedComponent.GetDependenciesTo(Component.Name)
		if e != nil {
			continue
		}
		for _, dependency := range dependencies {
			if e == nil {
				result += "\"" + relatedComponent.Name + "\" ->" + getGraphVizReference(dependency) + getEdgeLayoutFromDependency(dependency, relatedComponent.Display) + "\n"
			}
		}
	}

	// Draw infrastructure :-)
	for _, infrastructureDependency := range Component.InfrastructureDependencies {
		result = result + "\n\"" + infrastructureDependency.Type + "\"[shape=box, color=\"#576f96\"] \n"
		result = result + "\n\"" + infrastructureDependency.Type + "\"->\"" + Component.Name + "\"[color=\"#576f96\",arrowhead=none] \n"
	}
	result = result + "\n}"
	return result
}

func (ProjectDrawer *ProjectDrawer) drawComponentOutgoingRelations(Component *model.Application, hidePlanned bool) string {
	result := ""
	// Relation from components
	for _, dependency := range Component.Dependencies {
		if dependency.Status == model.STATUS_PLANNED && hidePlanned {
			continue
		}
		dependencyComponent, err := dependency.GetApplication(ProjectDrawer.originalProject)
		if err == nil && dependencyComponent.Status == model.STATUS_PLANNED && hidePlanned {
			continue
		}
		result += "\"" + Component.Name + "\" ->" + getGraphVizReference(dependency) + getEdgeLayoutFromDependency(dependency, Component.Display) + "\n"
	}
	// Relation from components/interfaces
	for _, providedInterface := range Component.ProvidedServices {
		for _, dependency := range providedInterface.Dependencies {
			if dependency.Status == model.STATUS_PLANNED && hidePlanned {
				continue
			}
			dependencyComponent, err := dependency.GetApplication(ProjectDrawer.originalProject)
			if err == nil && dependencyComponent.Status == model.STATUS_PLANNED && hidePlanned {
				continue
			}
			result += "\"" + Component.Name + "\":\"" + providedInterface.Name + "\"->" + getGraphVizReference(dependency) + getEdgeLayoutFromDependency(dependency, Component.Display) + "\n"
		}
	}
	return result
}

//Get GRaphviz style reference
func getGraphVizReference(Dependency model.Dependency) string {

	if strings.Contains(Dependency.Reference, ".") {
		splitted := strings.Split(Dependency.Reference, ".")
		return "\"" + splitted[0] + "\":\"" + splitted[1] + "\""
	}
	return "\"" + Dependency.Reference + "\""
}

func getEdgeLayoutFromDependency(dependency model.Dependency, display model.ApplicationDisplaySettings) string {

	edgeLayout := "["
	if display.BorderColor != "" {
		edgeLayout += "color=\"" + display.BorderColor + "\""
	} else {
		edgeLayout += "color=\"#333333\""
	}
	edgeLayout += ", fontsize=\"10\", fontcolor=\"#555555\" "
	if dependency.Relationship == "acl" {
		edgeLayout += ", dir=both, arrowtail=\"box\", taillabel=<<font color=\"red\" ><b>acl</b></font>>"
	} else {
		edgeLayout += ", label=\"" + dependency.Relationship + "\""
	}
	if dependency.Relationship == "customer-supplier" {
		edgeLayout += ", weight=2"
	}
	if dependency.Relationship == "conformist" || dependency.Relationship == "partnership" {
		edgeLayout += ", weight=3"
	}

	if dependency.Status == model.STATUS_PLANNED {
		edgeLayout += ", style=\"dotted\""
	} else if dependency.IsBrowserBased {
		edgeLayout += ", style=\"dashed\""
	}

	if dependency.IsSameLevel {
		edgeLayout += ", constraint=false"
	}
	return edgeLayout + "]"
}

// Factory
func CreateProjectDrawer(Project *model.Project, iconPath string) *ProjectDrawer {
	var Drawer ProjectDrawer
	Drawer.originalProject = Project
	Drawer.iconPath = iconPath
	return &Drawer
}
