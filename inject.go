package di

import (
	"reflect"
)

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

var injectableInterface = reflect.TypeOf(new(injectable)).Elem()

// canInject checks that typ is injectable
// Injectable type can be pointer to struct or struct and need to embed di.Inject.
func canInject(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		return typ.Implements(injectableInterface)
	}
	if typ.Kind() == reflect.Struct {
		return typ.Implements(injectableInterface)
	}
	return false
}

// parseFieldParams parses struct fields with di tag and form array of parameters associated
// with their field number.
func parseFieldParams(rt reflect.Type) (fields []int, params []parameter) {
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
	return fields, params
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
	return injectableStructFieldTag{
		name:     di,
		optional: optional == "true", // `optional:"true"`
	}, true
}

// createStructProvider creates embed provider.
func providerFromInjectableParameter(p parameter) *providerInjectable {
	var typ reflect.Type
	if p.typ.Kind() == reflect.Ptr {
		typ = p.typ.Elem()
	} else {
		typ = p.typ
	}
	provider := &providerInjectable{
		typ:  p.typ,
		name: p.name,
		val:  reflect.New(typ).Elem(),
	}
	provider.injectable.fields, provider.injectable.params = parseFieldParams(p.typ)
	return provider
}

type providerInjectable struct {
	typ        reflect.Type
	name       string
	injectable struct {
		// params parsed once
		params []parameter
		// field numbers parsed once
		fields []int
	}
	val reflect.Value
}

func (p providerInjectable) Type() reflect.Type {
	return p.typ
}

func (p providerInjectable) Name() string {
	return p.name
}

func (p providerInjectable) ParameterList() parameterList {
	return p.injectable.params
}

func (p providerInjectable) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	// set injectable fields
	if len(p.injectable.fields) > 0 {
		// result value
		rv := p.val
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if !rv.CanSet() {
			bug()
		}
		// field index
		for i, value := range values {
			// field value
			fv := rv.Field(p.injectable.fields[i])
			if !fv.CanSet() {
				continue
			}
			fv.Set(value)
		}
	}
	return p.val, nil, nil
}
