package di

import (
	"reflect"
)

// ctorType describes types of constructor provider.
type ctorType int

const (
	ctorUnknown           ctorType = iota
	ctorValue                      // (deps) (result)
	ctorValueError                 // (deps) (result, error)
	ctorValueCleanup               // (deps) (result, cleanup)
	ctorValueCleanupError          // (deps) (result, cleanup, error)
)

// constructorCompiler compiles constructor functions.
type constructorCompiler struct {
	typ ctorType
	fn  function
}

// newConstructorCompiler creates new function compiler from function.
func newConstructorCompiler(fn function) (*constructorCompiler, bool) {
	ctorType := determineCtorType(fn)
	if ctorType == ctorUnknown {
		return nil, false
	}
	return &constructorCompiler{
		typ: ctorType,
		fn:  fn,
	}, true
}

func (c constructorCompiler) deps(s schema) (deps []*node, err error) {
	for i := 0; i < c.fn.NumIn(); i++ {
		in := c.fn.Type.In(i)
		node, err := s.find(in, Tags{})
		if err != nil {
			return nil, err
		}
		deps = append(deps, node)
	}
	return deps, nil
}

func (c constructorCompiler) compile(dependencies []reflect.Value, s schema) (reflect.Value, error) {
	// call constructor function
	out := funcResult(c.fn.Call(dependencies))
	rv := out.value()
	switch c.typ {
	case ctorValue:
		return rv, nil
	case ctorValueError:
		return rv, out.error(1)
	case ctorValueCleanup:
		s.cleanup(out.cleanup())
		return rv, nil
	case ctorValueCleanupError:
		s.cleanup(out.cleanup())
		return rv, out.error(2)
	}
	bug()
	return reflect.Value{}, nil
}

func (c constructorCompiler) fields() map[int]field {
	return parsePopulateFields(c.fn.Out(0))
}

// determineCtorType
func determineCtorType(fn function) ctorType {
	switch true {
	case fn.NumOut() == 1:
		return ctorValue
	case fn.NumOut() == 2:
		if isError(fn.Out(1)) {
			return ctorValueError
		}
		if isCleanup(fn.Out(1)) {
			return ctorValueCleanup
		}
	case fn.NumOut() == 3 && isCleanup(fn.Out(1)) && isError(fn.Out(2)):
		return ctorValueCleanupError
	}
	return ctorUnknown
}

// funcResult is a helper struct for reflect.Call.
type funcResult []reflect.Value

// value returns first result type.
func (r funcResult) value() reflect.Value {
	return r[0]
}

// cleanup returns cleanup function.
func (r funcResult) cleanup() func() {
	if r[1].IsNil() {
		return nil
	}
	return r[1].Interface().(func())
}

// error returns error if it exists.
func (r funcResult) error(position int) error {
	if r[position].IsNil() {
		return nil
	}
	return r[position].Interface().(error)
}
