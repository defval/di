package di

import (
	"reflect"
)

type valueCompiler struct {
	rv reflect.Value
}

func (v valueCompiler) deps(s schema) ([]*node, error) {
	return nil, nil
}

func (v valueCompiler) compile(dependencies []reflect.Value, s schema) (reflect.Value, error) {
	return v.rv, nil
}
