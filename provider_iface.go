package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/reflection"
)

// providerInterface
type providerInterface struct {
	source keyUniq
	result key
}

// newProviderInterface
func newProviderInterface(uniq string, k key, as interface{}) (*providerInterface, error) {
	iface, err := reflection.InspectInterfacePtr(as)
	if err != nil {
		return nil, err
	}
	if !k.Type.Implements(iface.Type) {
		return nil, fmt.Errorf("%s not implement %s", k, iface.Type)
	}
	return &providerInterface{
		source: keyUniq{k, uniq},
		result: key{iface.Type, k.Name},
	}, nil
}

// Type
func (p providerInterface) Type() reflect.Type {
	return p.result.Type
}

// Name
func (p providerInterface) Name() string {
	return p.result.Name
}

// ParameterList
func (p providerInterface) ParameterList() parameterList {
	var plist parameterList
	plist = append(plist, parameter{
		uniq:     p.source.uniq,
		name:     p.source.Name,
		typ:      p.source.Type,
		optional: false,
	})
	return plist
}

// Provide
func (p providerInterface) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	return values[0], nil, nil
}
