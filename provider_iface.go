package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/reflection"
)

// providerInterface
type providerInterface struct {
	res      id
	provider provider
}

// newProviderInterface
func newProviderInterface(provider provider, as interface{}) (*providerInterface, error) {
	i, err := reflection.InspectInterfacePtr(as)
	if err != nil {
		return nil, err
	}
	if !provider.ID().Type.Implements(i.Type) {
		return nil, fmt.Errorf("%s not implement %s", provider.ID(), i.Type)
	}
	return &providerInterface{
		res: id{
			Name: provider.ID().Name,
			Type: i.Type,
		},
		provider: provider,
	}, nil
}

func (i *providerInterface) ID() id {
	return i.res
}

func (i *providerInterface) ParameterList() parameterList {
	var plist parameterList
	plist = append(plist, parameter{
		name:     i.provider.ID().Name,
		typ:      i.provider.ID().Type,
		optional: false,
		embed:    false,
	})
	return plist
}

func (i *providerInterface) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	return values[0], nil, nil
}
