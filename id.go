package di

import (
	"fmt"
	"reflect"
)

// key is a id of provider in container
type id struct {
	Name string
	Type reflect.Type
}

// String represent id as string.
func (i id) String() string {
	if i.Name == "" {
		return fmt.Sprintf("%s", i.Type)
	}
	return fmt.Sprintf("%s[%s]", i.Type, i.Name)
}
