package reflection

import (
	"fmt"
	"reflect"
)

// InspectInterfacePtr inspects interface pointer.
func InspectInterfacePtr(iface interface{}) (*Interface, error) {
	typ := reflect.TypeOf(iface)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Interface {
		return nil, fmt.Errorf("%s: not a pointer to interface", typ)
	}

	return &Interface{
		Name: typ.Elem().Name(),
		Type: typ.Elem(),
	}, nil
}

// Interface is a interface description.
type Interface struct {
	Name string
	Type reflect.Type
}
