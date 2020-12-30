package di

import (
	"fmt"
	"reflect"
)

// node is a dependency injection node.
type node struct {
	compiler
	rv     *reflect.Value
	rt     reflect.Type
	tags   Tags
	tracer Tracer
}

// parse parses fn constructor.
func nodeFromFunction(fn function) (*node, error) {
	cmp, ok := newConstructorCompiler(fn)
	if !ok {
		return nil, fmt.Errorf("invalid constructor signature, got %s", fn.Type)
	}
	// result type
	rt := fn.Out(0)
	// constructor result with di.Inject - only addressable pointers
	// anonymous parameters with di.Inject - only struct
	if isInjectable(rt) && rt.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("di.Inject not supported for unaddressable result of constructor, use *%s instead", rt)
	}
	tags := map[string]string{}
	if isTaggable(rt) {
		tt := rt
		if tt.Kind() == reflect.Ptr {
			tt = tt.Elem()
		}
		f, ok := tt.FieldByName("Tags")
		if !ok {
			return nil, fmt.Errorf("tags usage error: need to embed di.Tags without field name")
		}
		field, ok := inspectStructField(f)
		if ok {
			tags = field.tags
		}
	}
	return &node{
		rv:       &reflect.Value{},
		rt:       rt,
		tags:     tags,
		compiler: cmp,
	}, nil
}

// String is a string representation of node.
func (n *node) String() string {
	return fmt.Sprintf("%s%s", n.rt, n.tags)
}

// Build builds value of node.
func (n *node) Value(s schema) (reflect.Value, error) {
	if n.rv != nil && n.rv.IsValid() {
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
	*n.rv = rv
	if err := n.populate(s); err != nil {
		tracer.Trace("%s: %s", n.String(), err)
		return reflect.Value{}, err
	}
	tracer.Trace("Resolved %s", n.String())
	return rv, nil
}

// populate populates node fields.
func (n *node) populate(s schema) error {
	iv := *n.rv
	for i, field := range n.fields() {
		if iv.Kind() == reflect.Ptr {
			iv = iv.Elem()
		}
		node, err := s.find(field.rt, field.tags)
		if err != nil && field.optional {
			continue
		}
		if err != nil {
			return err
		}
		v, err := node.Value(s)
		if err != nil {
			return err
		}
		iv.Field(i).Set(v)
	}
	return nil
}

// getFunctionNodes
func getFunctionNodes(fn function, s schema) (params []*node, err error) {
	for i := 0; i < fn.NumIn(); i++ {
		in := fn.Type.In(i)
		node, err := s.find(in, Tags{})
		if err != nil {
			return nil, err
		}
		params = append(params, node)
	}
	return params, nil
}
