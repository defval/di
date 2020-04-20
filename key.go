package di

import (
	"fmt"
	"reflect"
)

// key is a type identity
type key struct {
	Type reflect.Type
	Name string
}

// String represent key as string.
func (i key) String() string {
	if i.Name == "" {
		return fmt.Sprintf("%s", i.Type)
	}
	return fmt.Sprintf("%s[%s]", i.Type, i.Name)
}

type keyUniq struct {
	key
	uniq string
}
