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

func (c typeCompiler) deps(s schema) (deps []*node, err error) {
	return nil, nil
}

func (c typeCompiler) compile(dependencies []reflect.Value, s schema) (reflect.Value, error) {
	return reflect.New(c.rt).Elem(), nil
}

func (c *typeCompiler) fields() map[int]field {
	return parseFields(c.rt)
}
