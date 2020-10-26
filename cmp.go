package di

import (
	"reflect"
)

// compiler compiles node result type.
type compiler interface {
	// params returns compiler params nodes.
	params(s schema) (params []*node, err error)
	// compile compiles resolved dependencies into value.
	compile(dependencies []reflect.Value, s schema) (reflect.Value, error)
}
