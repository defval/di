package di

import (
	"reflect"
)

type groupCompiler struct {
	rt      reflect.Type
	matched []*node
}

// newGroupCompiler creates group compiler of rt and with matched nodes.
func newGroupCompiler(rt reflect.Type, matched []*node) *groupCompiler {
	return &groupCompiler{
		rt:      rt,
		matched: matched,
	}
}

func (c *groupCompiler) params(s schema) (params []*node, err error) {
	return c.matched, nil
}

func (c *groupCompiler) compile(dependencies []reflect.Value, s schema) (reflect.Value, error) {
	return reflect.Append(reflect.New(c.rt).Elem(), dependencies...), nil
}

func (c *groupCompiler) fields() map[int]field {
	return nil
}
