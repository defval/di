package di

import (
	"fmt"
	"reflect"
	"runtime"
)

// Func is a function description.
type function struct {
	Name string
	reflect.Type
	reflect.Value
}

var errorInterface = reflect.TypeOf(new(error)).Elem()

// isError checks that typ have error signature.
func isError(typ reflect.Type) bool {
	return typ.Implements(errorInterface)
}

// isCleanup checks that typ have cleanup signature.
func isCleanup(typ reflect.Type) bool {
	return typ.Kind() == reflect.Func && typ.NumIn() == 0 && typ.NumOut() == 0
}

// InspectFunc inspects function.
func inspectFunction(fn interface{}) (function, bool) {
	if reflect.ValueOf(fn).Kind() != reflect.Func {
		return function{}, false
	}
	val := reflect.ValueOf(fn)
	typ := val.Type()
	funcForPC := runtime.FuncForPC(val.Pointer())
	return function{
		Name:  funcForPC.Name(),
		Type:  typ,
		Value: val,
	}, true
}

// Interface is a interface description.
type link struct {
	Name string
	Type reflect.Type
}

// inspectInterfacePointer inspects interface pointer.
func inspectInterfacePointer(i interface{}) (*link, error) {
	if i == nil {
		return nil, fmt.Errorf("nil: not a pointer to interface")
	}
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Interface {
		return nil, fmt.Errorf("%s: not a pointer to interface", typ)
	}

	return &link{
		Name: typ.Elem().Name(),
		Type: typ.Elem(),
	}, nil
}
