package reflection

import (
	"reflect"
	"runtime"
)

// IsFunc check that value have function type.
func IsFunc(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Func
}

// Func is a function description.
type Func struct {
	Name string
	reflect.Type
	reflect.Value
}

// InspectFunc inspects function.
func InspectFunc(fn interface{}) (Func, bool) {
	if !IsFunc(fn) {
		return Func{}, false
	}
	val := reflect.ValueOf(fn)
	typ := val.Type()
	funcForPC := runtime.FuncForPC(val.Pointer())
	return Func{
		Name:  funcForPC.Name(),
		Type:  typ,
		Value: val,
	}, true
}
