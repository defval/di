package di

import (
	"fmt"
	"reflect"
)

// providerStub
type providerStub struct {
	msg string
	res key
}

// newProviderStub
func newProviderStub(k key, msg string) *providerStub {
	return &providerStub{res: k, msg: msg}
}

func (m *providerStub) Key() key {
	return m.res
}

func (m *providerStub) ParameterList() parameterList {
	return parameterList{}
}

func (m *providerStub) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	return reflect.Value{}, nil, fmt.Errorf(m.msg)
}
