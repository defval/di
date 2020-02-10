package reflection

import (
	"fmt"
	"reflect"
	"runtime"
)

// IsFunc
func IsFunc(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Func
}

// Func
type Func struct {
	Name string
	reflect.Type
	reflect.Value
}

// InspectFunction
func InspectFunction(fn interface{}) *Func {
	if !IsFunc(fn) {
		panic(fmt.Sprintf("%s: not a function", reflect.TypeOf(fn).Kind())) // todo: improve message
	}

	val := reflect.ValueOf(fn)
	fnpc := runtime.FuncForPC(val.Pointer())

	return &Func{
		Name:  fnpc.Name(),
		Type:  val.Type(),
		Value: val,
	}
}
