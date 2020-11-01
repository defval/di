package di

import (
	"reflect"
)

type typeCompiler struct {
	rt reflect.Type
}

// newTypeCompiler creates compiler that creates new instance of rt.
func newTypeCompiler(rt reflect.Type) *typeCompiler {
	return &typeCompiler{rt: rt}
}

func (p typeCompiler) params(s schema) (params []*node, err error) {
	return nil, nil
}

func (p typeCompiler) compile(dependencies []reflect.Value, s schema) (reflect.Value, error) {
	return reflect.New(p.rt).Elem(), nil
}
