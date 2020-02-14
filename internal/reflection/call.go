package reflection

import "reflect"

// CallResult is a helper struct for reflect.Call.
type CallResult []reflect.Value

// Result returns first result type.
func (r CallResult) Result() reflect.Value {
	return r[0]
}

// Cleanup returns cleanup function.
func (r CallResult) Cleanup() func() {
	if r[1].IsNil() {
		return nil
	}
	return r[1].Interface().(func())
}

// Error returns error if it exists.
func (r CallResult) Error(position int) error {
	if r[position].IsNil() {
		return nil
	}
	return r[position].Interface().(error)
}
