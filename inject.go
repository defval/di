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
	return injectableStructFieldTag{
		name:     di,
		optional: optional == "true", // `optional:"true"`
	}, true
}

// createStructProvider creates embed provider.
func providerFromInjectableParameter(p parameter) *providerEmbed {
	var typ reflect.Type
	if p.typ.Kind() == reflect.Ptr {
		typ = p.typ.Elem()
	} else {
		typ = p.typ
	}
	provider := &providerEmbed{
		id: id{
			Name: p.name,
			Type: p.typ,
		},
		typ: typ,
		val: reflect.New(typ).Elem(),
	}
	provider.injectable.params, provider.injectable.fields = parseInjectableType(p.typ)
	return provider
}

type providerEmbed struct {
	id         id
	typ        reflect.Type
	val        reflect.Value
	injectable struct {
		// params parsed once
		params []parameter
		// field numbers parsed once
		fields []int
	}
}

func (p *providerEmbed) ID() id {
	return p.id
}

func (p *providerEmbed) ParameterList() parameterList {
	return p.injectable.params
}

func (p *providerEmbed) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	// set injectable fields
	if len(p.injectable.fields) > 0 {
		// result value
		rv := p.val
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if !rv.CanSet() {
			panic("you found a bug, please create new issue for this: https://github.com/goava/di/issues/new")
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
