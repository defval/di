package graphkv

import (
	"github.com/emicklei/dot"
)

// Node
type Node struct {
	Key   Key
	Value interface{}
}

// Graph
type Graph struct {
	dag    *directedGraph
	values map[Key]interface{}
}

// New
func New() *Graph {
	return &Graph{
		dag:    newDirectedGraph(),
		values: map[Key]interface{}{},
	}
}

// Get
func (g *Graph) Get(key Key) Node {
	return Node{Key: key, Value: g.values[key]}
}

// Replace
func (g *Graph) Replace(key Key, value interface{}) {
	g.values[key] = value
}

// Add
func (g *Graph) Add(key Key, value interface{}) {
	g.dag.AddNode(key)
	g.values[key] = value
}

// Edge
func (g *Graph) Edge(from Key, to Key) {
	g.dag.AddEdge(from, to)
}

// Exists
func (g *Graph) Exists(key Key) bool {
	return g.dag.NodeExists(key)
}

// Nodes
func (g *Graph) Nodes() []Node {
	var nodes []Node
	for _, key := range g.dag.Nodes() {
		nodes = append(nodes, Node{key, g.values[key]})
	}
	return nodes
}

// CheckCycles
func (g *Graph) CheckCycles() error {
	_, err := g.dag.DFSSort()
	return err // todo: errors
}

// DOTGraph
func (g *Graph) DOTGraph() *dot.Graph {
	return g.dag.DOTGraph()
}
