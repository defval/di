package reflection

import "reflect"

var errorInterface = reflect.TypeOf(new(error)).Elem()

// IsError checks that typ have error signature.
func IsError(typ reflect.Type) bool {
	return typ.Implements(errorInterface)
}

// IsCleanup checks that typ have cleanup signature.
func IsCleanup(typ reflect.Type) bool {
	return typ.Kind() == reflect.Func && typ.NumIn() == 0 && typ.NumOut() == 0
}
