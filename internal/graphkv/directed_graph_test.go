package graphkv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestDirectedGraph() *directedGraph {
	graph := newDirectedGraph()
	graph.AddNodes("A", "B", "C", "D")
	return graph
}

func TestNewDirectedGraph(t *testing.T) {
	graph := newDirectedGraph()
	assert.NotNil(t, graph, "graph should not be nil")
	assert.Zero(t, graph.NodeCount(), "graph.NodeCount() should equal zero")
	assert.Empty(t, graph.Nodes(), "graph.Nodes() should equal empty")
	assert.Zero(t, graph.EdgeCount(), "graph.EdgeCount() should equal zero")
}

func TestDirectedGraphAddEdge(t *testing.T) {
	graph := newTestDirectedGraph()
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "D")
	graph.AddEdge("C", "B")

	assert.Equal(t, 3, graph.EdgeCount(), "graph.EdgeCount() should equal 3")
	assert.True(t, graph.EdgeExists("A", "B"), "graph.EdgeExists(A, B) should equal true")
	assert.True(t, graph.EdgeExists("B", "D"), "graph.EdgeExists(B, D) should equal true")
	assert.True(t, graph.EdgeExists("C", "B"), "graph.EdgeExists(C, B) should equal true")
}

func TestDirectedGraphAddEdgeDuplicate(t *testing.T) {
	graph := newTestDirectedGraph()
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "C")
	graph.AddEdge("B", "C")

	assert.Equal(t, 2, graph.EdgeCount(), "graph.EdgeCount() should equal 2")
	assert.True(t, graph.EdgeExists("A", "B"), "graph.EdgeExists(A, B) should equal true")
	assert.True(t, graph.EdgeExists("B", "C"), "graph.EdgeExists(B, C) should equal true")
}

func TestDirectedGraphAddEdgeMissingNodes(t *testing.T) {
	graph := newDirectedGraph()
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "C")

	assert.Equal(t, 3, graph.NodeCount(), "graph.NodeCount() should equal 2")
	assert.Equal(t, 2, graph.EdgeCount(), "graph.EdgeCount() should equal 2")
	assert.True(t, graph.EdgeExists("A", "B"), "graph.EdgeExists(A, B) should equal true")
	assert.True(t, graph.EdgeExists("B", "C"), "graph.EdgeExists(B, C) should equal true")
}

func TestDirectedGraphRemoveEdge(t *testing.T) {
	graph := newTestDirectedGraph()
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "D")
	graph.AddEdge("C", "B")
	graph.RemoveEdge("A", "B")
	graph.RemoveEdge("C", "B")

	assert.Equal(t, 1, graph.EdgeCount(), "graph.EdgeCount() should equal 1")
	assert.False(t, graph.EdgeExists("A", "B"), "graph.EdgeExists(A, B) should equal false")
	assert.False(t, graph.EdgeExists("C", "B"), "graph.EdgeExists(C, B) should equal false")
	assert.True(t, graph.EdgeExists("B", "D"), "graph.EdgeExists(B, D) should equal true")
}

func TestDirectedGraphRemoveEdgeMissing(t *testing.T) {
	graph := newTestDirectedGraph()
	graph.AddEdge("C", "B")
	graph.RemoveEdge("D", "A")
	graph.RemoveEdge("C", "B")

	assert.Zero(t, graph.EdgeCount(), "graph.EdgeCount() should equal zero")
}

func TestDirectedGraphHasEdges(t *testing.T) {
	graph := newTestDirectedGraph()
	graph.AddEdge("A", "C")

	assert.True(t, graph.HasEdges("A"), "graph.HasEdges(A) should equal true")
	assert.False(t, graph.HasEdges("B"), "graph.HasEdges(B) should equal false")
	assert.True(t, graph.HasEdges("C"), "graph.HasEdges(C) should equal true")
	assert.False(t, graph.HasEdges("D"), "graph.HasEdges(D) should equal false")
}

func TestDirectedGraphIncomingEdgeCount(t *testing.T) {
	graph := newTestDirectedGraph()
	graph.AddEdge("A", "C")
	graph.AddEdge("B", "C")

	assert.Zero(t, graph.IncomingEdgeCount("A"), "graph.IncomingEdgeCount(A) should equal 0")
	assert.Zero(t, graph.IncomingEdgeCount("B"), "graph.IncomingEdgeCount(B) should equal 0")
	assert.Equal(t, 2, graph.IncomingEdgeCount("C"), "graph.IncomingEdgeCount(C) should equal 1")
	assert.Zero(t, graph.IncomingEdgeCount("D"), "graph.IncomingEdgeCount(D) should equal 0")
}

func TestDirectedGraphOutgoingEdgeCount(t *testing.T) {
	graph := newTestDirectedGraph()
	graph.AddEdge("A", "B")
	graph.AddEdge("A", "C")

	assert.Equal(t, 2, graph.OutgoingEdgeCount("A"), "graph.OutgoingEdgeCount(A) should equal 2")
	assert.Zero(t, graph.OutgoingEdgeCount("B"), "graph.OutgoingEdgeCount(B) should equal 0")
	assert.Zero(t, graph.OutgoingEdgeCount("C"), "graph.OutgoingEdgeCount(C) should equal 0")
	assert.Zero(t, graph.OutgoingEdgeCount("D"), "graph.OutgoingEdgeCount(D) should equal 0")
}

func TestDirectedGraphRootNodes(t *testing.T) {
	graph := newTestDirectedGraph()
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "C")
	graph.AddEdge("D", "C")
	graph.AddEdge("E", "C")
	graph.AddEdge("F", "E")

	assert.Equal(t, []Key{"A", "D", "F"}, graph.RootNodes(), "graph.RootNodes() should equal [A, D, F]")
}

func TestDirectedGraphIsolatedNodes(t *testing.T) {
	graph := newTestDirectedGraph()
	graph.AddEdge("A", "C")

	assert.Equal(t, []Key{"B", "D"}, graph.IsolatedNodes(), "graph.IsolatedNodes() should equal [B, D]")
}

func TestDirectedGraphAdjacencyMatrix(t *testing.T) {
	graph := newTestDirectedGraph()
	graph.AddEdge("A", "C")
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "D")
	graph.AddEdge("C", "A")
	graph.AddEdge("D", "D")

	expected := map[interface{}]map[interface{}]bool{
		"A": map[interface{}]bool{"A": false, "B": true, "C": true, "D": false},
		"B": map[interface{}]bool{"A": false, "B": false, "C": false, "D": true},
		"C": map[interface{}]bool{"A": true, "B": false, "C": false, "D": false},
		"D": map[interface{}]bool{"D": true, "A": false, "B": false, "C": false},
	}

	assert.Equal(t, expected, graph.AdjacencyMatrix(), "graph.AdjacencyMatrix() should equal [B, D]")
}
