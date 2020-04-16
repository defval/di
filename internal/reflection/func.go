package reflection

import (
	"reflect"
	"runtime"
)

// Func is a function description.
type Func struct {
	Name string
	reflect.Type
	reflect.Value
}

// InspectFunc inspects function.
func InspectFunc(fn interface{}) (Func, bool) {
	if reflect.ValueOf(fn).Kind() != reflect.Func {
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
