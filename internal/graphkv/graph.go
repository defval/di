package graphkv

type graph struct {
	nodes *nodeList
}

func newGraph() *graph {
	return &graph{
		nodes: newNodeList(),
	}
}

// Copy returns a clone of the graph.
func (g *graph) Copy() *graph {
	return &graph{
		nodes: g.nodes.Copy(),
	}
}

// Nodes returns the graph's nodes.
// The slice is mutable for performance reasons but should not be mutated.
func (g *graph) Nodes() []Key {
	return g.nodes.Nodes()
}

// NodeCount returns the number of nodes.
func (g *graph) NodeCount() int {
	return g.nodes.Count()
}

// AddNode inserts the specified node into the graph.
// A node can be any value, e.g. int, string, pointer to a struct, map etc.
// Duplicate nodes are ignored.
func (g *graph) AddNode(node Key) {
	g.AddNodes(node)
}

// AddNodes inserts the specified nodes into the graph.
// A node can be any value, e.g. int, string, pointer to a struct, map etc.
// Duplicate nodes are ignored.
func (g *graph) AddNodes(nodes ...Key) {
	g.nodes.Add(nodes...)
}

// RemoveNode removes the specified nodes from the graph.
// If the node does not exist within the graph the call will fail silently.
func (g *graph) RemoveNode(node Key) {
	g.RemoveNodes(node)
}

// RemoveNodes removes the specified nodes from the graph.
// If a node does not exist within the graph the call will fail silently.
func (g *graph) RemoveNodes(nodes ...Key) {
	g.nodes.Remove(nodes...)
}

// NodeExists determines whether the specified node exists within the graph.
func (g *graph) NodeExists(node Key) bool {
	return g.nodes.Exists(node)
}
