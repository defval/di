package graph

import (
	"os"
	"testing"

	"github.com/goava/di/internal/graph/testgraph"
)

func TestGraph_CheckCycles(t *testing.T) {
	for _, graph := range testgraph.GraphSlice {
		f, err := os.Open("testdata/graph.json")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()
		g, err := NewGraphFromJSON(f, graph.Name)
		if err != nil {
			t.Error(err)
		}
		isDAG := true
		if err := CheckCycles(g); err != nil {
			isDAG = false
		}
		if isDAG != graph.IsDAG {
			t.Errorf("%s | IsDag are supposed to be %v", graph.Name, graph.IsDAG)
		}
	}
}
