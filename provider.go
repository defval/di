package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/graph"
)

// provider is a internal provider interface.
type provider interface {
	ID() id
	// ParameterList returns array of dependencies.
	ParameterList() parameterList
	// Provide provides value from provided parameters.
	Provide(values ...reflect.Value) (_ reflect.Value, cleanup func(), err error)
}

// providerNode is a adapter for graph node.
// todo: remove this, refactor graph to work with provider directly
type providerNode struct {
	provider
}

// ID returns graph ID
func (n providerNode) ID() graph.ID {
	return n.provider.ID()
}

// String returns graph provider string representation.
func (n providerNode) String() string {
	return fmt.Sprintf("%s provider", n.ID())
}
