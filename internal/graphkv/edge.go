package graphkv

type directedEdgeList struct {
	outgoingEdges map[Key]*nodeList
	incomingEdges map[Key]*nodeList
}

func newDirectedEdgeList() *directedEdgeList {
	return &directedEdgeList{
		outgoingEdges: make(map[Key]*nodeList),
		incomingEdges: make(map[Key]*nodeList),
	}
}

func (l *directedEdgeList) Copy() *directedEdgeList {
	outgoingEdges := make(map[Key]*nodeList, len(l.outgoingEdges))
	for node, edges := range l.outgoingEdges {
		outgoingEdges[node] = edges.Copy()
	}

	incomingEdges := make(map[Key]*nodeList, len(l.incomingEdges))
	for node, edges := range l.incomingEdges {
		incomingEdges[node] = edges.Copy()
	}

	return &directedEdgeList{
		outgoingEdges: outgoingEdges,
		incomingEdges: incomingEdges,
	}
}

func (l *directedEdgeList) Count() int {
	return len(l.outgoingEdges)
}

func (l *directedEdgeList) HasOutgoingEdges(node Key) bool {
	_, ok := l.outgoingEdges[node]
	return ok
}

func (l *directedEdgeList) OutgoingEdgeCount(node Key) int {
	if list := l.outgoingNodeList(node, false); list != nil {
		return list.Count()
	}
	return 0
}

func (l *directedEdgeList) outgoingNodeList(node Key, create bool) *nodeList {
	if list, ok := l.outgoingEdges[node]; ok {
		return list
	}
	if create {
		list := newNodeList()
		l.outgoingEdges[node] = list
		return list
	}
	return nil
}

func (l *directedEdgeList) OutgoingEdges(node Key) []Key {
	if list := l.outgoingNodeList(node, false); list != nil {
		return list.Nodes()
	}
	return nil
}

func (l *directedEdgeList) HasIncomingEdges(node Key) bool {
	_, ok := l.incomingEdges[node]
	return ok
}

func (l *directedEdgeList) IncomingEdgeCount(node Key) int {
	if list := l.incomingNodeList(node, false); list != nil {
		return list.Count()
	}
	return 0
}

func (l *directedEdgeList) incomingNodeList(node Key, create bool) *nodeList {
	if list, ok := l.incomingEdges[node]; ok {
		return list
	}
	if create {
		list := newNodeList()
		l.incomingEdges[node] = list
		return list
	}
	return nil
}

func (l *directedEdgeList) IncomingEdges(node Key) []Key {
	if list := l.incomingNodeList(node, false); list != nil {
		return list.Nodes()
	}
	return nil
}

func (l *directedEdgeList) Add(from Key, to Key) {
	l.outgoingNodeList(from, true).Add(to)
	l.incomingNodeList(to, true).Add(from)
}

func (l *directedEdgeList) Remove(from Key, to Key) {
	if list := l.outgoingNodeList(from, false); list != nil {
		list.Remove(to)

		if list.Count() == 0 {
			delete(l.outgoingEdges, from)
		}
	}
	if list := l.incomingNodeList(to, false); list != nil {
		list.Remove(from)

		if list.Count() == 0 {
			delete(l.incomingEdges, to)
		}
	}
}

func (l *directedEdgeList) Exists(from Key, to Key) bool {
	if list := l.outgoingNodeList(from, false); list != nil {
		return list.Exists(to)
	}
	return false
}
