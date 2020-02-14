package di

import (
	"fmt"
	"reflect"
)

// providerStub
type providerStub struct {
	id  id
	msg string
}

// newProviderStub
func newProviderStub(id id, msg string) *providerStub {
	return &providerStub{id: id, msg: msg}
}

func (m *providerStub) ID() id {
	return m.id
}

func (m *providerStub) ParameterList() parameterList {
	return parameterList{}
}

func (m *providerStub) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	return reflect.Value{}, nil, fmt.Errorf(m.msg)
}
