package graphkv

import (
	"errors"
)

// Errors relating to the DFSSorter.
var (
	ErrCyclicGraph = errors.New("the graph cannot be cyclic")
)

// DFSSorter topologically sorts a directed graph's nodes based on the
// directed edges between them using the Depth-first search algorithm.
type DFSSorter struct {
	graph      *directedGraph
	sorted     []Key
	visiting   map[Key]bool
	discovered map[Key]bool
}

// NewDFSSorter returns a new DFS sorter.
func NewDFSSorter(graph *directedGraph) *DFSSorter {
	return &DFSSorter{
		graph: graph,
	}
}

func (s *DFSSorter) init() {
	s.sorted = make([]Key, 0, s.graph.NodeCount())
	s.visiting = make(map[Key]bool)
	s.discovered = make(map[Key]bool, s.graph.NodeCount())
}

// Sort returns the sorted nodes.
func (s *DFSSorter) Sort() ([]Key, error) {
	s.init()

	// > while there are unmarked nodes do
	for _, node := range s.graph.Nodes() {
		if err := s.visit(node); err != nil {
			return nil, err
		}
	}

	// as the nodes were appended to the slice for performance reasons,
	// rather than prepended as correctly stated by the algorithm,
	// we need to reverse the sorted slice
	for i, j := 0, len(s.sorted)-1; i < j; i, j = i+1, j-1 {
		s.sorted[i], s.sorted[j] = s.sorted[j], s.sorted[i]
	}

	return s.sorted, nil
}

// See https://en.wikipedia.org/wiki/Topological_sorting#Depth-first_search
func (s *DFSSorter) visit(node Key) error {
	// > if n has a permanent mark then return
	if discovered, ok := s.discovered[node]; ok && discovered {
		return nil
	}
	// > if n has a temporary mark then stop (not a DAG)
	if visiting, ok := s.visiting[node]; ok && visiting {
		return ErrCyclicGraph
	}

	// > mark n temporarily
	s.visiting[node] = true

	// > for each node m with an edge from n to m do
	for _, outgoing := range s.graph.OutgoingEdges(node) {
		if err := s.visit(outgoing); err != nil {
			return err
		}
	}

	s.discovered[node] = true
	delete(s.visiting, node)

	s.sorted = append(s.sorted, node)
	return nil
}

// DFSSort returns the graph's nodes in topological order based on the
// directed edges between them using the Depth-first search algorithm.
func (g *directedGraph) DFSSort() ([]Key, error) {
	sorter := NewDFSSorter(g)
	return sorter.Sort()
}

// Errors relating to the CoffmanGrahamSorter.
var (
	ErrDependencyOrder = errors.New("the topological dependency order is incorrect")
)

// CoffmanGrahamSorter sorts a graph's nodes into a sequence of levels,
// arranging so that a node which comes after another in the order is
// assigned to a lower level, and that a level never exceeds the width.
// See https://en.wikipedia.org/wiki/Coffman–Graham_algorithm
type CoffmanGrahamSorter struct {
	graph *directedGraph
	width int
}

// NewCoffmanGrahamSorter returns a new Coffman-Graham sorter.
func NewCoffmanGrahamSorter(graph *directedGraph, width int) *CoffmanGrahamSorter {
	return &CoffmanGrahamSorter{
		graph: graph,
		width: width,
	}
}

// Sort returns the sorted nodes.
func (s *CoffmanGrahamSorter) Sort() ([][]Key, error) {
	// create a copy of the graph and remove transitive edges
	reduced := s.graph.Copy()
	reduced.RemoveTransitives()

	// topologically sort the graph nodes
	nodes, err := reduced.DFSSort()
	if err != nil {
		return nil, err
	}

	layers := make([][]Key, 0)
	levels := make(map[Key]int, len(nodes))

	for _, node := range nodes {
		dependantLevel := -1
		for _, dependant := range reduced.IncomingEdges(node) {
			level, ok := levels[dependant]
			if !ok {
				return nil, ErrDependencyOrder
			}
			if level > dependantLevel {
				dependantLevel = level
			}
		}

		level := -1
		// find the first unfilled layer outgoing the dependent layer
		// skip this if the dependent layer is the last
		if dependantLevel < len(layers)-1 {
			for i := dependantLevel + 1; i < len(layers); i++ {
				// ensure the layer doesn't exceed the desired width
				if len(layers[i]) < s.width {
					level = i
					break
				}
			}
		}
		// create a new layer new none was found
		if level == -1 {
			layers = append(layers, make([]Key, 0, 1))
			level = len(layers) - 1
		}

		layers[level] = append(layers[level], node)
		levels[node] = level
	}

	return layers, nil
}

// CoffmanGrahamSort sorts the graph's nodes into a sequence of levels,
// arranging so that a node which comes after another in the order is
// assigned to a lower level, and that a level never exceeds the specified width.
func (g *directedGraph) CoffmanGrahamSort(width int) ([][]Key, error) {
	sorter := NewCoffmanGrahamSorter(g, width)
	return sorter.Sort()
}
