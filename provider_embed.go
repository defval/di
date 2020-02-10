package di

import (
	"reflect"
	"strings"
)

// createStructProvider
func newProviderEmbed(p parameter) *providerEmbed {
	var embedType reflect.Type
	if p.res.Kind() == reflect.Ptr {
		embedType = p.res.Elem()
	} else {
		embedType = p.res
	}

	return &providerEmbed{
		key: key{
			name: p.name,
			res:  p.res,
			typ:  ptEmbedParameter,
		},
		embedType:  embedType,
		embedValue: reflect.New(embedType).Elem(),
	}
}

type providerEmbed struct {
	key        key
	embedType  reflect.Type
	embedValue reflect.Value
}

func (p *providerEmbed) Key() key {
	return p.key
}

func (p *providerEmbed) ParameterList() parameterList {
	var plist parameterList
	for i := 0; i < p.embedType.NumField(); i++ {
		name, optional, isDependency := p.inspectFieldTag(i)
		if !isDependency {
			continue
		}
		field := p.embedType.Field(i)
		plist = append(plist, parameter{
			name:     name,
			res:      field.Type,
			optional: optional,
			embed:    isEmbedParameter(field.Type),
		})
	}
	return plist
}

func (p *providerEmbed) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	for i, offset := 0, 0; i < p.embedType.NumField(); i++ {
		_, _, isDependency := p.inspectFieldTag(i)
		if !isDependency {
			offset++
			continue
		}

		p.embedValue.Field(i).Set(values[i-offset])
	}

	return p.embedValue, nil, nil
}

func (p *providerEmbed) inspectFieldTag(num int) (name string, optional bool, isDependency bool) {
	fieldType := p.embedType.Field(num)
	fieldValue := p.embedValue.Field(num)
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
