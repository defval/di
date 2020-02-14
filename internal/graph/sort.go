package graph

// TopologicalSort does topological sort(ordering) with DFS.
// It returns true if the graph is a DAG (no cycle, with a topological sort).
// False if the graph is not a DAG (cycle, with no topological sort).
func TopologicalSort(g *Graph) ([]ID, bool) {
	// L = Empty list that will contain the sorted nodes
	var L []ID
	isDAG := true
	color := make(map[ID]string)
	for v := range g.Nodes() {
		color[v] = "white"
	}

	// for each vertex v in G:
	for v := range g.Nodes() {
		// if v.color == "white":
		if color[v] == "white" {
			// topologicalSortVisit(v, L, isDAG)
			topologicalSortVisit(g, v, &L, &isDAG, &color)
		}
	}

	return L, isDAG
}

func topologicalSortVisit(
	g *Graph,
	id ID,
	L *[]ID,
	isDAG *bool,
	color *map[ID]string,
) {
	// if v.color == "gray":
	if (*color)[id] == "gray" {
		// isDAG = false
		*isDAG = false
		return
	}

	// if v.color == "white":
	if (*color)[id] == "white" {
		// v.color = "gray":
		(*color)[id] = "gray"

		// for each child vertex w of v:
		cmap, err := g.Targets(id)
		if err != nil {
			panic(err)
		}
		for w := range cmap {
			topologicalSortVisit(g, w, L, isDAG, color)
		}

		// v.color = "black"
		(*color)[id] = "black"

		// L.push_front(v)
		temp := make([]ID, len(*L)+1)
		temp[0] = id
		copy(temp[1:], *L)
		*L = temp
	}
}
