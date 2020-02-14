package di

import (
	"reflect"
	"strings"
)

// createStructProvider creates embed provider.
func providerFromEmbedParameter(p parameter) *providerEmbed {
	var embedType reflect.Type
	if p.typ.Kind() == reflect.Ptr {
		embedType = p.typ.Elem()
	} else {
		embedType = p.typ
	}

	return &providerEmbed{
		id: id{
			Name: p.name,
			Type: p.typ,
		},
		typ: embedType,
		val: reflect.New(embedType).Elem(),
	}
}

type providerEmbed struct {
	id  id
	typ reflect.Type
	val reflect.Value
}

func (p *providerEmbed) ID() id {
	return p.id
}

func (p *providerEmbed) ParameterList() parameterList {
	var plist parameterList
	for i := 0; i < p.typ.NumField(); i++ {
		name, optional, isDependency := p.inspectFieldTag(i)
		if !isDependency {
			continue
		}
		field := p.typ.Field(i)
		plist = append(plist, parameter{
			name:     name,
			typ:      field.Type,
			optional: optional,
			embed:    isEmbedParameter(field.Type),
		})
	}
	return plist
}

func (p *providerEmbed) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	for i, offset := 0, 0; i < p.typ.NumField(); i++ {
		_, _, isDependency := p.inspectFieldTag(i)
		if !isDependency {
			offset++
			continue
		}

		p.val.Field(i).Set(values[i-offset])
	}

	return p.val, nil, nil
}

func (p *providerEmbed) inspectFieldTag(num int) (name string, optional bool, isDependency bool) {
	fieldType := p.typ.Field(num)
	fieldValue := p.val.Field(num)
	tag, tagExists := fieldType.Tag.Lookup("di")
	if !tagExists || !fieldValue.CanSet() {
		return "", false, false
	}
	name, optional = p.parseTag(tag)
	return name, optional, true
}

func (p *providerEmbed) parseTag(tag string) (name string, optional bool) {
	options := strings.Split(tag, ",")
	if len(options) == 0 {
		return "", false
	}
	if len(options) == 1 && options[0] == "optional" {
		return "", true
	}
	if len(options) == 1 {
		return options[0], false
	}
	if len(options) == 2 && options[1] == "optional" {
		return options[0], true
	}
	panic("incorrect di tag")
}
