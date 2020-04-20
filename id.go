package di

import (
	"fmt"
	"reflect"
)

// id is a type identity
type id struct {
	Type reflect.Type
	Name string
}

// String represent id as string.
func (i id) String() string {
	if i.Name == "" {
		return fmt.Sprintf("%s", i.Type)
	}
	return fmt.Sprintf("%s[%s]", i.Type, i.Name)
}
