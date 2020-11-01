package di

import (
	"reflect"
	"strings"
)

var iInjectable = reflect.TypeOf(new(injectable)).Elem()

// Inject indicates that public fields with special tag type will be injected automatically.
//
//	type MyType struct {
//		di.Inject
//
//		Server *http.Server `di:""` // will be injected
//	}
type Inject struct {
	injectable
}

// injectable interface needs to struct fields injection functional.
type injectable interface {
	isInjectable()
}

// isInjectable checks that typ is injectable
// Injectable type can be pointer to struct or struct and need to embed di.Inject.
func isInjectable(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		return typ.Implements(iInjectable)
	}
	if typ.Kind() == reflect.Struct {
		return typ.Implements(iInjectable)
	}
	return false
}

type field struct {
	rt       reflect.Type
	tags     Tags
	optional bool
}

// `name:"asd" command:"console"`
func parseField(f reflect.StructField) (field, bool) {
	skip, _ := f.Tag.Lookup("skip")
	if skip == "true" {
		return field{}, false
	}
	tag := string(f.Tag)
	kvs := strings.Split(tag, " ")
	tags := Tags{}
	if len(kvs) == 0 {
		return field{}, false
	}
	for _, v := range kvs {
		kv := strings.Split(v, ":")
		if len(kv) != 2 {
			continue
		}
		k := kv[0]
		v := strings.Trim(kv[1], "\"")
		tags[k] = v
	}
	optional, _ := f.Tag.Lookup("optional")
	return field{
		rt:       f.Type,
		tags:     tags,
		optional: optional == "true",
	}, true
}
