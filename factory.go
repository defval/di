package di

import "reflect"

// validateOptionsFactory validates function.
func validateOptionsFactory(fn function) bool {
	if fn.NumOut() == 1 && isOptionsSlice(fn.Out(0)) {
		return true
	}
	if fn.NumOut() == 2 && isOptionsSlice(fn.Out(0)) && isError(fn.Out(1)) {
		return true
	}
	return false
}

func isOptionsSlice(typ reflect.Type) bool {
	optType := reflect.TypeOf((*Option)(nil)).Elem()
	return typ.Kind() == reflect.Slice &&
		typ.Elem().Kind() == reflect.Interface &&
		typ.Elem().Implements(optType)
}
