package di

import (
	"reflect"
)

// newProviderGroup creates new group from provided resultKey.
func newProviderGroup(k key) *providerGroup {
	ifaceKey := key{
		res: reflect.SliceOf(k.res),
		typ: ptGroup,
	}

	return &providerGroup{
		result: ifaceKey,
		pl:     parameterList{},
	}
}

// providerGroup
type providerGroup struct {
	result key
	pl     parameterList
}

// Add
func (i *providerGroup) Add(k key) {
	i.pl = append(i.pl, parameter{
		name:     k.name,
		res:      k.res,
		optional: false,
		embed:    false,
	})
}

// resultKey
func (i providerGroup) Key() key {
	return i.result
}

// parameters
func (i providerGroup) ParameterList() parameterList {
	return i.pl
}

// Provide
func (i providerGroup) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	group := reflect.New(i.result.res).Elem()
	return reflect.Append(group, values...), nil, nil
}
