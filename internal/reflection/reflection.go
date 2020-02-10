package reflection

import "reflect"

var errorInterface = reflect.TypeOf(new(error)).Elem()

// IsError
func IsError(typ reflect.Type) bool {
	return typ.Implements(errorInterface)
}

// IsCleanup
func IsCleanup(typ reflect.Type) bool {
	return typ.Kind() == reflect.Func && typ.NumIn() == 0 && typ.NumOut() == 0
}

// IsPtr
func IsPtr(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Ptr
}
