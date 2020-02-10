package graphkv

// directedGraph is a graph supporting directed edges between nodes.
type directedGraph struct {
	*graph
	edges *directedEdgeList
}

// newDirectedGraph creates a graph of nodes with directed edges.
func newDirectedGraph() *directedGraph {
	return &directedGraph{
		graph: newGraph(),
		edges: newDirectedEdgeList(),
	}
}

// Copy returns a clone of the directed graph.
func (g *directedGraph) Copy() *directedGraph {
	return &directedGraph{
		graph: g.graph.Copy(),
		edges: g.edges.Copy(),
	}
}

// EdgeCount returns the number of direced edges between nodes.
func (g *directedGraph) EdgeCount() int {
	return g.edges.Count()
}

// AddEdge adds the edge to the graph.
func (g *directedGraph) AddEdge(from Key, to Key) {
	// prevent adding an edge referring to missing nodes
	if !g.NodeExists(from) {
		g.AddNode(from)
	}
	if !g.NodeExists(to) {
		g.AddNode(to)
	}

	g.edges.Add(from, to)
}

// RemoveEdge removes the edge from the graph.
func (g *directedGraph) RemoveEdge(from Key, to Key) {
	g.edges.Remove(from, to)
}

// HasEdges determines whether the graph contains any edges to or from the node.
func (g *directedGraph) HasEdges(node Key) bool {
	if g.HasIncomingEdges(node) {
		return true
	}
	return g.HasOutgoingEdges(node)
}

// EdgeExists checks whether the edge exists within the graph.
func (g *directedGraph) EdgeExists(from Key, to Key) bool {
	return g.edges.Exists(from, to)
}

// HasIncomingEdges checks whether the graph contains any directed
// edges pointing to the node.
func (g *directedGraph) HasIncomingEdges(node Key) bool {
	return g.edges.HasIncomingEdges(node)
}

// IncomingEdges returns the nodes belonging to directed edges pointing
// towards the specified node.
func (g *directedGraph) IncomingEdges(node Key) []Key {
	return g.edges.IncomingEdges(node)
}

// IncomingEdgeCount returns the number of edges pointing from the specified
// node (indegree).
func (g *directedGraph) IncomingEdgeCount(node Key) int {
	return g.edges.IncomingEdgeCount(node)
}

// HasOutgoingEdges checks whether the graph contains any directed
// edges pointing from the node.
func (g *directedGraph) HasOutgoingEdges(node Key) bool {
	return g.edges.HasOutgoingEdges(node)
}

// OutgoingEdges returns the nodes belonging to directed edges pointing
// from the specified node.
func (g *directedGraph) OutgoingEdges(node Key) []Key {
	return g.edges.OutgoingEdges(node)
}

// OutgoingEdgeCount returns the number of edges pointing from the specified
// node (outdegree).
func (g *directedGraph) OutgoingEdgeCount(node Key) int {
	return g.edges.OutgoingEdgeCount(node)
}

// RootNodes finds the entry-point nodes to the graph, i.e. those without
// incoming edges.
func (g *directedGraph) RootNodes() []Key {
	results := make([]Key, 0)
	for _, node := range g.Nodes() {
		if !g.HasIncomingEdges(node) {
			results = append(results, node)
		}
	}
	return results
}

// IsolatedNodes finds independent nodes in the graph, i.e. those without edges.
func (g *directedGraph) IsolatedNodes() []Key {
	results := make([]Key, 0)
	for _, node := range g.Nodes() {
		if !g.HasEdges(node) {
			results = append(results, node)
		}
	}
	return results
}

// AdjacencyMatrix returns a matrix indicating whether pairs of nodes are
// adjacent or not within the graph.
func (g *directedGraph) AdjacencyMatrix() map[Key]map[Key]bool {
	matrix := make(map[Key]map[Key]bool, g.NodeCount())
	for _, a := range g.Nodes() {
		matrix[a] = make(map[Key]bool, g.NodeCount())

		for _, b := range g.Nodes() {
			matrix[a][b] = g.EdgeExists(a, b)
		}
	}
	return matrix
}

// RemoveTransitives removes any transitive edges so that as fewest possible
// edges exist while matching the reachability of the original graph.
func (g *directedGraph) RemoveTransitives() {
	for _, a := range g.Nodes() {
		for _, b := range g.Nodes() {
			if !g.EdgeExists(a, b) {
				continue
			}
			for _, c := range g.Nodes() {
				if g.EdgeExists(b, c) {
					g.RemoveEdge(a, c)
				}
			}
		}
	}
}
