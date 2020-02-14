package di

import (
	"reflect"
)

// newProviderGroup creates new group from provided key.
func newProviderGroup(k id) *providerGroup {
	id := id{
		Type: reflect.SliceOf(k.Type), // creates []<type> group
	}
	return &providerGroup{
		id: id,
		pl: parameterList{},
	}
}

// providerGroup
type providerGroup struct {
	id id
	pl parameterList
}

func (i *providerGroup) ID() id {
	return i.id
}

// Add
func (i *providerGroup) Add(k id) {
	i.pl = append(i.pl, parameter{
		name:     k.Name,
		typ:      k.Type,
		optional: false,
		embed:    false,
	})
}

// parameters
func (i providerGroup) ParameterList() parameterList {
	return i.pl
}

// Provide
func (i providerGroup) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	group := reflect.New(i.id.Type).Elem()
	return reflect.Append(group, values...), nil, nil
}
