package graphviz

import (
	"testing"

	"strings"

	"github.com/AOEpeople/vistecture/v2/model/core"
)

func TestApplicationDrawer_Draw(t *testing.T) {

	project := core.Project{
		Name: "Project1",
		Applications: []*core.Application{
			{
				Name: "app1",

				Dependencies: []core.Dependency{
					{
						Reference: "app2",
					},
				},
			},
			{
				Name: "app2",
			},
			{
				Name: "app3",
			},
		},
	}

	drawer := CreateProjectDrawer(&project, "")
	if graph := drawer.DrawComplete(false); !strings.Contains(graph, "graph [") {
		t.Error("Graph contains no graph [] declaration ", graph)
	}

	if graph := drawer.DrawComplete(false); !strings.Contains(graph, "\"app1\" ->\"app2\"") {
		t.Error("Graph contains no edge", graph)
	}
	if graph := drawer.DrawComplete(false); !strings.Contains(graph, "\"app3\"") {
		t.Error("Graph contains no core app3", graph)
	}
}
