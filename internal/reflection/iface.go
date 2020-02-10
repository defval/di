package reflection

import (
	"fmt"
	"reflect"
)

// InspectInterfacePtr
func InspectInterfacePtr(iface interface{}) *Interface {
	typ := reflect.TypeOf(iface)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Interface {
		panic(fmt.Sprintf("%s: not a pointer to interface", typ)) // todo: improve message
	}

	return &Interface{
		Name: typ.Elem().Name(),
		Type: typ.Elem(),
	}
}

// Interface
type Interface struct {
	Name string
	Type reflect.Type
}
