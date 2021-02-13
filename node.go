package di

import (
	"fmt"
	"reflect"
)

// newConstructorNode
func newConstructorNode(ctor interface{}) (*node, error) {
	f, valid := inspectFunction(ctor)
	if !valid {
		return nil, fmt.Errorf("invalid constructor signature, got %s", reflect.TypeOf(ctor))
	}
	cmp, ok := newConstructorCompiler(f)
	if !ok {
		return nil, fmt.Errorf("invalid constructor signature, got %s", f.Type)
	}
	// result type
	rt := f.Out(0)
	if _, err := canInject(rt); err != nil {
		return nil, err
	}
	tags := map[string]string{}
	if haveTags(rt) {
		tmp := rt
		if tmp.Kind() == reflect.Ptr {
			tmp = tmp.Elem()
		}
		f, ok := tmp.FieldByName("Tags")
		if !ok {
			return nil, fmt.Errorf("tags usage error: need to embed di.Tags without field name")
		}
		field, ok := inspectStructField(f)
		if ok {
			tags = field.tags
		}
	}
	return &node{
		rv:       new(reflect.Value),
		rt:       rt,
		tags:     tags,
		compiler: cmp,
	}, nil
}

// node is a dependency injection node.
type node struct {
	compiler
	rt   reflect.Type
	tags Tags
	// rv value can be shared between nodes
	// initializing node always need to allocate memory for rv
	rv *reflect.Value
}

// String is a string representation of node.
func (n *node) String() string {
	return fmt.Sprintf("%s%s", n.rt, n.tags)
}

// Build builds value of node.
func (n *node) Value(s schema) (reflect.Value, error) {
	tracer.Trace("-- %s requested", n)
	if n.rv.IsValid() {
		tracer.Trace("-- %s already compiled", n)
		return *n.rv, nil
	}
	nodes, err := n.deps(s)
	if err != nil {
		tracer.Trace("%s: %s", n.String(), err)
		return reflect.Value{}, err
	}
	var dependencies []reflect.Value
	for _, node := range nodes {
		v, err := node.Value(s)
		if err != nil {
			tracer.Trace("%s: %s", n.String(), err)
			return reflect.Value{}, err
		}
		dependencies = append(dependencies, v)
	}
	rv, err := n.compile(dependencies, s)
	if err != nil {
		tracer.Trace("%s: %s", n.String(), err)
		return reflect.Value{}, err
	}
	// if result value not addr, create pointer for it
	if !rv.CanAddr() {
		addr := reflect.New(rv.Type())
		addr.Elem().Set(rv)
		rv = addr.Elem()
	}
	*n.rv = rv
	if err := populate(s, *n.rv); err != nil {
		tracer.Trace("%s: %s", n.String(), err)
		return reflect.Value{}, err
	}
	tracer.Trace("Resolved %s", n.String())
	return *n.rv, nil
}

// populate populates node fields.
func populate(s schema, rv reflect.Value) error {
	if !rv.IsValid() {
		panic("node zero result value on populate")
	}
	inject, err := canInject(rv.Type())
	if err != nil {
		return err
	}
	if !inject {
		return nil
	}
	// indirect pointer
	if rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}
	for index, field := range parsePopulateFields(rv.Type()) {
		node, err := s.find(field.rt, field.tags)
		if err != nil && field.optional {
			tracer.Trace("-- Skip optional field: %s", field)
			continue
		}
		if err != nil {
			return err
		}
		v, err := node.Value(s)
		if err != nil {
			return err
		}
		f := rv.Field(index)
		if !f.CanSet() {
			panic(fmt.Sprintf("can not set field %s(%d) of %s (addr: %t)", f.Type(), f.Pointer(), rv.Type(), rv.CanAddr()))
		}
		f.Set(v)
	}
	return nil
}
