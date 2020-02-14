package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/goava/di/internal/graph/testgraph"
)

// NewGraphFromJSON returns a new Graph from a JSON file.
// Here's the sample JSON data:
//
//	{
//	    "graph_00": {
//	        "S": {
//	            "A": 100,
//	            "B": 14,
//	            "C": 200
//	        },
//	        "A": {
//	            "S": 15,
//	            "B": 5,
//	            "D": 20,
//	            "T": 44
//	        },
//	        "B": {
//	            "S": 14,
//	            "A": 5,
//	            "D": 30,
//	            "E": 18
//	        },
//	        "C": {
//	            "S": 9,
//	            "E": 24
//	        },
//	        "D": {
//	            "A": 20,
//	            "B": 30,
//	            "E": 2,
//	            "F": 11,
//	            "T": 16
//	        },
//	        "E": {
//	            "B": 18,
//	            "C": 24,
//	            "D": 2,
//	            "F": 6,
//	            "T": 19
//	        },
//	        "F": {
//	            "D": 11,
//	            "E": 6,
//	            "T": 6
//	        },
//	        "T": {
//	            "A": 44,
//	            "D": 16,
//	            "F": 6,
//	            "E": 19
//	        }
//	    },
//	}
//
func NewGraphFromJSON(rd io.Reader, graphID string) (*Graph, error) {
	js := make(map[string]map[string]map[string]float64)
	dec := json.NewDecoder(rd)
	for {
		if err := dec.Decode(&js); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	if _, ok := js[graphID]; !ok {
		return nil, fmt.Errorf("%s does not exist", graphID)
	}
	gmap := js[graphID]

	g := New()
	for id1, mm := range gmap {
		nd1, err := g.Node(StringID(id1))
		if err != nil {
			nd1 = NewNode(id1)
			g.AddNode(nd1)
		}
		for id2, weight := range mm {
			nd2, err := g.Node(StringID(id2))
			if err != nil {
				nd2 = NewNode(id2)
				g.AddNode(nd2)
			}
			g.ReplaceEdge(nd1.ID(), nd2.ID(), weight)
		}
	}

	return g, nil
}

type node struct {
	id string
}

func NewNode(id string) Node {
	return &node{
		id: id,
	}
}

func (n *node) ID() ID {
	return StringID(n.id)
}

func (n *node) String() string {
	return n.id
}

func TestNewGraph(t *testing.T) {
	g1 := New()
	fmt.Println("g1:", g1.String())

	f, err := os.Open("testdata/graph.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	g2, err := NewGraphFromJSON(f, "graph_00")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("g2:", g2.String())
}

func TestNewGraphFromJSON_graph(t *testing.T) {
	f, err := os.Open("testdata/graph.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	g, err := NewGraphFromJSON(f, "graph_00")
	if err != nil {
		t.Fatalf("nil graph %v", err)
	}
	if g.nodeToTargets[StringID("C")][StringID("S")] != 9.0 {
		t.Fatalf("weight from C to S must be 9.0 but %f", g.nodeToTargets[StringID("C")][StringID("S")])
	}
	for _, tg := range testgraph.GraphSlice {
		f, err := os.Open("testdata/graph.json")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		g, err := NewGraphFromJSON(f, tg.Name)
		if err != nil {
			t.Fatalf("nil graph %v", err)
		}
		if g.NodeCount() != tg.TotalNodeCount {
			t.Fatalf("%s | Expected %d but %d", tg.Name, tg.TotalNodeCount, g.NodeCount())
		}
		for _, elem := range tg.EdgeToWeight {
			weight1, err := g.Weight(StringID(elem.Nodes[0]), StringID(elem.Nodes[1]))
			if err != nil {
				t.Fatal(err)
			}
			weight2 := elem.Weight
			if weight1 != weight2 {
				t.Fatalf("Expected %f but %f", weight2, weight1)
			}
		}
	}
}

func TestGraph_GetVertices(t *testing.T) {
	for _, tg := range testgraph.GraphSlice {
		f, err := os.Open("testdata/graph.json")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		g, err := NewGraphFromJSON(f, tg.Name)
		if err != nil {
			t.Fatal(err)
		}
		if g.NodeCount() != tg.TotalNodeCount {
			t.Fatalf("wrong number of vertices: %s", g)
		}
	}
}

func TestGraph_DeleteNode(t *testing.T) {
	f, err := os.Open("testdata/graph.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	g, err := NewGraphFromJSON(f, "graph_01")
	if err != nil {
		t.Fatal(err)
	}
	if !g.DeleteNode(StringID("D")) {
		t.Fatal("D does not exist in the graph")
	}
	if nd, err := g.Node(StringID("D")); err == nil {
		t.Fatalf("No Node Expected but got %s", nd)
	}
	if v, err := g.Sources(StringID("C")); err != nil || len(v) != 1 {
		t.Fatalf("Expected 1 edge incoming to C but %v\n\n%s", err, g)
	}
	if v, err := g.Targets(StringID("C")); err != nil || len(v) != 2 {
		t.Fatalf("Expected 2 edges outgoing from C but %v\n\n%s", err, g)
	}
	if v, err := g.Targets(StringID("F")); err != nil || len(v) != 2 {
		t.Fatalf("Expected 2 edges outgoing from F but %v\n\n%s", err, g)
	}
	if v, err := g.Sources(StringID("F")); err != nil || len(v) != 2 {
		t.Fatalf("Expected 2 edges incoming to F but %v\n\n%s", err, g)
	}
	if v, err := g.Targets(StringID("B")); err != nil || len(v) != 3 {
		t.Fatalf("Expected 3 edges outgoing from B but %v\n\n%s", err, g)
	}
	if v, err := g.Sources(StringID("E")); err != nil || len(v) != 4 {
		t.Fatalf("Expected 4 edges incoming to E but %v\n\n%s", err, g)
	}
	if v, err := g.Targets(StringID("E")); err != nil || len(v) != 3 {
		t.Fatalf("Expected 3 edges outgoing from E but %v\n\n%s", err, g)
	}
	if v, err := g.Targets(StringID("T")); err != nil || len(v) != 3 {
		t.Fatalf("Expected 3 edges outgoing from T but %v\n\n%s", err, g)
	}
}

func TestGraph_DeleteEdge(t *testing.T) {
	f, err := os.Open("testdata/graph.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	g, err := NewGraphFromJSON(f, "graph_01")
	if err != nil {
		t.Fatal(err)
	}

	if err := g.DeleteEdge(StringID("B"), StringID("D")); err != nil {
		t.Fatal(err)
	}
	if v, err := g.Sources(StringID("D")); err != nil || len(v) != 4 {
		t.Fatalf("Expected 4 edges incoming to D but %v\n\n%s", err, g)
	}

	if err := g.DeleteEdge(StringID("B"), StringID("C")); err != nil {
		t.Fatal(err)
	}
	if err := g.DeleteEdge(StringID("S"), StringID("C")); err != nil {
		t.Fatal(err)
	}
	if v, err := g.Targets(StringID("S")); err != nil || len(v) != 2 {
		t.Fatalf("Expected 2 edges outgoing from S but %v\n\n%s", err, g)
	}

	if err := g.DeleteEdge(StringID("C"), StringID("E")); err != nil {
		t.Fatal(err)
	}
	if err := g.DeleteEdge(StringID("E"), StringID("D")); err != nil {
		t.Fatal(err)
	}
	if v, err := g.Targets(StringID("E")); err != nil || len(v) != 3 {
		t.Fatalf("Expected 3 edges outgoing from E but %v\n\n%s", err, g)
	}
	if v, err := g.Sources(StringID("E")); err != nil || len(v) != 3 {
		t.Fatalf("Expected 3 edges incoming to E but %v\n\n%s", err, g)
	}

	if err := g.DeleteEdge(StringID("F"), StringID("E")); err != nil {
		t.Fatal(err)
	}
	if v, err := g.Sources(StringID("E")); err != nil || len(v) != 2 {
		t.Fatalf("Expected 2 edges incoming to E but %v\n\n%s", err, g)
	}
}

func TestGraph_ReplaceEdge(t *testing.T) {
	f, err := os.Open("testdata/graph.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	g, err := NewGraphFromJSON(f, "graph_00")
	if err != nil {
		t.Fatal(err)
	}
	if err := g.ReplaceEdge(StringID("C"), StringID("S"), 1.0); err != nil {
		t.Fatal(err)
	}
	if v, err := g.Weight(StringID("C"), StringID("S")); err != nil || v != 1.0 {
		t.Fatalf("weight from C to S must be 1.0 but %v\n\n%v", err, g)
	}
}
