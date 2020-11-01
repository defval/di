package di

import (
	"fmt"
	"reflect"
)

// node is a dependency injection node.
type node struct {
	rv   *reflect.Value
	rt   reflect.Type
	tags Tags
	compiler
}

// parse parses fn constructor.
func nodeFromFunction(fn function) (*node, error) {
	cmp, ok := newFuncCompiler(fn)
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
		field, ok := parseField(f)
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
	nodes, err := n.params(s)
	if err != nil {
		return reflect.Value{}, err
	}
	var dependencies []reflect.Value
	for _, node := range nodes {
		v, err := node.Value(s)
		if err != nil {
			return reflect.Value{}, err
		}
		dependencies = append(dependencies, v)
	}
	rv, err := n.compile(dependencies, s)
	if err != nil {
		return reflect.Value{}, err
	}
	*n.rv = rv
	if err := n.populate(s); err != nil {
		return reflect.Value{}, err
	}
	return rv, nil
}

func (n *node) fields() map[int]field {
	rt := n.rt
	if !isInjectable(rt) {
		return nil
	}
	rv := *n.rv
	if !rv.IsValid() {
		switch rt.Kind() {
		case reflect.Ptr:
			rv = reflect.New(rt.Elem())
		default:
			rv = reflect.New(rt).Elem()
		}

	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	fields := make(map[int]field, rt.NumField())
	// fi - field index
	for fi := 0; fi < rt.NumField(); fi++ {
		fv := rv.Field(fi)
		// check that field can be set
		if !fv.CanSet() || rt.Field(fi).Anonymous {
			continue
		}
		// cur - current field
		cur := rt.Field(fi)
		f, valid := parseField(cur)
		if !valid {
			continue
		}
		fields[fi] = field{
			rt:       cur.Type,
			tags:     f.tags,
			optional: f.optional,
		}
	}
	return fields
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
