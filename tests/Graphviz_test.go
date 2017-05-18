package tests

import (
	"strings"
	"testing"

	"github.com/AOEpeople/vistecture/model/core"
	"github.com/AOEpeople/vistecture/model/graphviz"
)

func TestCanGetGraph(t *testing.T) {

	project, e := core.CreateProject("fixture")
	if e != nil {
		t.Error("Factory returned error", e)
	}
	drawer := graphviz.CreateProjectDrawer(project, "")
	if graph := drawer.DrawComplete(); !strings.Contains(graph, "graph [") {
		t.Error("Graph contains no graph [] declaration ", graph)
	}

	if graph := drawer.DrawComplete(); !strings.Contains(graph, "\"app1\" ->\"app2\"") {
		t.Error("Graph contains no edge", graph)
	}
}
