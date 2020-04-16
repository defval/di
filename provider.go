package di

import (
	"reflect"
)

// provider is a internal provider interface.
type provider interface {
	ID() id
	// ParameterList returns array of dependencies.
	ParameterList() parameterList
	// Provide provides value from provided parameters.
	Provide(values ...reflect.Value) (_ reflect.Value, cleanup func(), err error)
}
