package di

import (
	"reflect"
)

// compiler compiles dependency node.
type compiler interface {
	// deps return array of nodes that will be used for node compilation.
	deps(s schema) ([]*node, error)
	// compile compiles node. The dependencies are already compiled dependencies of this type.
	compile(dependencies []reflect.Value, s schema) (reflect.Value, error)
}
