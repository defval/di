package di

import (
	"io"

	"github.com/emicklei/dot"
)

// Graph
type Graph struct {
	graph *dot.Graph
}

func (g *Graph) WriteTo(writer io.Writer) {
	g.graph.Write(writer)
}

func (g *Graph) String() string {
	return g.graph.String()
}
