package graphkv

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGraph(t *testing.T) {
	graph := newGraph()
	assert.NotNil(t, graph, "graph should not be nil")
	assert.Zero(t, graph.NodeCount(), "graph.NodeCount() should equal zero")
	assert.Empty(t, graph.Nodes(), "graph.Nodes() should equal empty")
}

func TestGraphAddNode(t *testing.T) {
	graph := newGraph()
	graph.AddNode("A")
	graph.AddNode("B")
	graph.AddNode("C")

	assert.Equal(t, 3, graph.NodeCount(), "graph.NodeCount() should equal 3")
	assert.Equal(t, []interface{}{"A", "B", "C"}, graph.Nodes(), "graph.Nodes() should equal [A, B, C]")
}

func TestGraphAddNodeDuplicate(t *testing.T) {
	graph := newGraph()
	graph.AddNode("A")
	graph.AddNode("B")
	graph.AddNode("C")
	graph.AddNode("A")

	assert.Equal(t, 3, graph.NodeCount(), "graph.NodeCount() should equal 3")
	assert.Equal(t, []interface{}{"A", "B", "C"}, graph.Nodes(), "graph.Nodes() should equal [A, B, C]")
}

func BenchmarkGraphAddNodes(b *testing.B) {
	for i := 12.0; i <= 20; i++ {
		count := int(math.Pow(2, i))

		b.Run(fmt.Sprintf("%d", count), func(b *testing.B) {
			graph := newGraph()
			for i := 0; i < count; i++ {
				graph.AddNode(i)
			}
		})
	}
}

func TestGraphAddNodes(t *testing.T) {
	graph := newGraph()
	graph.AddNodes("A", "B", "C")

	assert.Equal(t, 3, graph.NodeCount(), "graph.NodeCount() should equal 3")
	assert.Equal(t, []interface{}{"A", "B", "C"}, graph.Nodes(), "graph.Nodes() should equal [A, B, C]")
}

func TestGraphRemoveNode(t *testing.T) {
	graph := newGraph()
	graph.AddNode("A")
	graph.AddNode("B")
	graph.AddNode("C")
	graph.AddNode("D")
	graph.RemoveNode("A")
	graph.RemoveNode("C")

	assert.Equal(t, 2, graph.NodeCount(), "graph.NodeCount() should equal 2")
	assert.Equal(t, []interface{}{"B", "D"}, graph.Nodes(), "graph.Nodes() should equal [B, D]")
}

func TestGraphRemoveNodeMissing(t *testing.T) {
	graph := newGraph()
	graph.AddNode("A")
	graph.AddNode("B")
	graph.AddNode("C")
	graph.AddNode("D")
	graph.RemoveNode("A")
	graph.RemoveNode("A")
	graph.RemoveNode("E")

	assert.Equal(t, 3, graph.NodeCount(), "graph.NodeCount() should equal 2")
	assert.Equal(t, []interface{}{"B", "C", "D"}, graph.Nodes(), "graph.Nodes() should equal [B, C, D]")
}

func TestGraphRemoveNodes(t *testing.T) {
	graph := newGraph()
	graph.AddNode("A")
	graph.AddNode("B")
	graph.AddNode("C")
	graph.AddNode("D")
	graph.RemoveNodes("A", "C")

	assert.Equal(t, 2, graph.NodeCount(), "graph.NodeCount() should equal 2")
	assert.Equal(t, []interface{}{"B", "D"}, graph.Nodes(), "graph.Nodes() should equal [B, D]")
}

func TestGraphNodeExists(t *testing.T) {
	graph := newGraph()
	assert.False(t, graph.NodeExists("A"), "graph.NodeExists(\"A\") should equal false")
	assert.False(t, graph.NodeExists("B"), "graph.NodeExists(\"B\") should equal false")

	graph.AddNode("A")
	assert.True(t, graph.NodeExists("A"), "graph.NodeExists(\"A\") should equal true")
	assert.False(t, graph.NodeExists("B"), "graph.NodeExists(\"B\") should equal false")

	graph.RemoveNode("A")
	assert.False(t, graph.NodeExists("A"), "graph.NodeExists(\"A\") should equal false")
	assert.False(t, graph.NodeExists("B"), "graph.NodeExists(\"B\") should equal false")
}
