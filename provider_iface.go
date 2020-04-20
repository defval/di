package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/reflection"
)

// providerInterface
type providerInterface struct {
	typ      reflect.Type
	name     string
	provider provider
}

// newProviderInterface
func newProviderInterface(provider provider, as interface{}) (*providerInterface, error) {
	i, err := reflection.InspectInterfacePtr(as)
	if err != nil {
		return nil, err
	}
	if !provider.Type().Implements(i.Type) {
		return nil, fmt.Errorf("%s not implement %s", provider.Type(), i.Type)
	}
	return &providerInterface{
		typ:      i.Type,
		name:     provider.Name(),
		provider: provider,
	}, nil
}

func (i providerInterface) Type() reflect.Type {
	return i.typ
}

func (i providerInterface) Name() string {
	return i.name
}

func (i providerInterface) ParameterList() parameterList {
	var plist parameterList
	plist = append(plist, parameter{
		name:     i.provider.Name(),
		typ:      i.provider.Type(),
		optional: false,
	})
	return plist
}

func (i providerInterface) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	return values[0], nil, nil
}
