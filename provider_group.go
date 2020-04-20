package di

import (
	"reflect"
)

// newProviderGroup creates new group from provided key.
func newProviderGroup(gtype reflect.Type, plist *providerList) *providerGroup {
	var params parameterList
	for _, p := range plist.Uniqs() {
		params = append(params, parameter{
			uniq:     p.uniq,
			typ:      p.Type(),
			name:     p.Name(),
			optional: false,
		})
	}
	return &providerGroup{
		typ:    gtype,
		params: params,
	}
}

// providerGroup
type providerGroup struct {
	typ    reflect.Type
	params parameterList
}

func (p providerGroup) Type() reflect.Type {
	return p.typ
}

func (p providerGroup) Name() string {
	return ""
}

func (p providerGroup) ParameterList() parameterList {
	return p.params
}

func (p providerGroup) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	group := reflect.New(p.typ).Elem()
	return reflect.Append(group, values...), nil, nil
}
