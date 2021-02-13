package di

import (
	"reflect"
)

type typeCompiler struct {
	rt reflect.Type
}

// newTypeCompiler creates compiler that creates new instance of rt.
func newTypeCompiler(rt reflect.Type) *typeCompiler {
	return &typeCompiler{
		rt: rt,
	}
}

func (c typeCompiler) deps(s schema) (deps []*node, err error) {
	return nil, nil
}

func (c typeCompiler) compile(dependencies []reflect.Value, s schema) (reflect.Value, error) {
	if c.rt.Kind() == reflect.Ptr {
		rt := c.rt.Elem()
		zero := reflect.Zero(rt)
		addr := reflect.New(rt)
		addr.Elem().Set(zero)
		return addr, nil

	}
	return reflect.New(c.rt).Elem(), nil
}

func (c *typeCompiler) fields() map[int]field {
	return parsePopulateFields(c.rt)
}
