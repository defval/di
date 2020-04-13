package di

import (
	"reflect"
	"strings"
)

// Parameter is embed helper that indicates that type is a constructor embed parameter.
// Deprecated: Use di.Injectable
type Parameter struct {
	injectable
}

// Injectable indicates that public fields with special tag type will be injected automatically.
//
//	type MyType struct {
//		di.Injectable
//
//		Server *http.Server `di:""` // will be injected
//	}
type Injectable struct {
	injectable
}

// internalParameter
type injectable interface {
	isInjectable()
}

var injectableInterface = reflect.TypeOf(new(injectable)).Elem()

// isInjectable checks that typ is injectable.
func isInjectable(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		return typ.Implements(injectableInterface)
	}
	if typ.Kind() == reflect.Struct {
		return typ.Implements(injectableInterface)
	}
	return false
}

func parseInjectableType(rt reflect.Type) (params []parameter, fields []int) {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	// create instance for check private fields.
	rv := reflect.New(rt).Elem()
	// fi - field index
	for fi := 0; fi < rt.NumField(); fi++ {
		fv := rv.Field(fi)
		// check that field can be set
		if !fv.CanSet() || rt.Field(fi).Anonymous {
			continue
		}
		// cur - current field
		cur := rt.Field(fi)
		parsed, valid := inspectInjectableStructFieldType(cur)
		if !valid {
			continue
		}
		params = append(params, parameter{
			name:     parsed.name,
			typ:      cur.Type,
			optional: parsed.optional,
		})
		fields = append(fields, fi)
	}
	return params, fields
}

// injectableStructFieldTag contains injectable field params
type injectableStructFieldTag struct {
	name     string // `name:"my_dep"`
	optional bool   // `optional:"true"`
}

// inspectInjectableStructFieldType inspects injectable struct field and parse tags.
func inspectInjectableStructFieldType(field reflect.StructField) (injectableStructFieldTag, bool) {
	di, exists := field.Tag.Lookup("di")
	if !exists {
		return injectableStructFieldTag{}, false
	}
	optional, _ := field.Tag.Lookup("optional")
	options := strings.Split(di, ",")
	if len(options) == 0 {
		return injectableStructFieldTag{
			optional: optional == "true",
		}, true
	}
	if len(options) == 1 && options[0] == "optional" {
		return injectableStructFieldTag{
			optional: true,
		}, true
	}
	if len(options) == 1 {
		return injectableStructFieldTag{
			name:     options[0],
			optional: optional == "true",
		}, true
	}
	if len(options) == 2 && options[1] == "optional" {
		return injectableStructFieldTag{
			name:     options[0],
			optional: true,
		}, true
	}
	return injectableStructFieldTag{
		name:     di,
		optional: optional == "true", // `optional:"true"`
	}, true
}
