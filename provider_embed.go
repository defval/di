package di

import (
	"reflect"
)

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
