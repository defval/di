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

func (g *groupCompiler) params(s schema) (params []*node, err error) {
	return g.matched, nil
}

func (g *groupCompiler) compile(dependencies []reflect.Value, s schema) (reflect.Value, error) {
	return reflect.Append(reflect.New(g.rt).Elem(), dependencies...), nil
}
