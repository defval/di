package graphkv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDFSSorter(t *testing.T) {
	graph := newDirectedGraph()
	graph.AddNodes(0, 1, 2, 3, 4, 5, 6, 7)
	graph.AddEdge(0, 2)
	graph.AddEdge(1, 2)
	graph.AddEdge(1, 5)
	graph.AddEdge(1, 6)
	graph.AddEdge(2, 5)
	graph.AddEdge(3, 5)
	graph.AddEdge(5, 6)
	graph.AddEdge(5, 7)

	sorted, err := graph.DFSSort()

	assert.NoError(t, err, "graph.DFSSort() error should be nil")
	assert.Equal(t, []Key{4, 3, 1, 0, 2, 5, 7, 6}, sorted, "graph.DFSSort() nodes should equal [4, 3, 1, 0, 2, 5, 7, 6]")
}

func TestDFSSorterCyclic(t *testing.T) {
	graph := newDirectedGraph()
	graph.AddNodes(0, 1)
	graph.AddEdge(0, 1)
	graph.AddEdge(1, 0)

	sorted, err := graph.DFSSort()

	assert.EqualError(t, err, ErrCyclicGraph.Error(), "graph.DFSSort() error should be ErrCyclicGraph")
	assert.Nil(t, sorted, "graph.DFSSort() nodes should be nil")
}

func TestCoffmanGrahamSorter(t *testing.T) {
	graph := newDirectedGraph()

	graph.AddNodes(0, 1, 2, 3, 4, 5, 6, 7, 8)
	graph.AddEdge(0, 2)
	graph.AddEdge(0, 5)
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)
	graph.AddEdge(2, 4)
	graph.AddEdge(3, 6)
	graph.AddEdge(4, 6)
	graph.AddEdge(5, 7)
	graph.AddEdge(6, 7)
	graph.AddEdge(6, 8)

	sorted, err := graph.CoffmanGrahamSort(2)

	assert.NoError(t, err, "graph.CoffmanGrahamSort(2)0 error should be nil")
	assert.Equal(t, [][]Key{
		[]Key{1, 0},
		[]Key{5, 2},
		[]Key{4, 3},
		[]Key{6},
		[]Key{8, 7},
	}, sorted, "graph.CoffmanGrahamSort(2) nodes should equal [[1, 0], [5, 2], [4, 3], [6], [8, 7]]")
}

func TestCoffmanGrahamSorterCyclic(t *testing.T) {
	graph := newDirectedGraph()

	graph.AddNodes(0, 1, 2, 3, 4, 5, 6, 7, 8)
	graph.AddEdge(0, 2)
	graph.AddEdge(0, 5)
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 0) // cyclic edge
	graph.AddEdge(2, 3)
	graph.AddEdge(2, 4)
	graph.AddEdge(3, 6)
	graph.AddEdge(4, 6)
	graph.AddEdge(5, 7)
	graph.AddEdge(6, 7)
	graph.AddEdge(6, 8)

	sorted, err := graph.CoffmanGrahamSort(2)

	assert.EqualError(t, err, ErrCyclicGraph.Error(), "graph.CoffmanGrahamSort(2) error should be ErrCyclicGraph")
	assert.Nil(t, sorted, "graph.CoffmanGrahamSort(2) nodes should be nil")
}
