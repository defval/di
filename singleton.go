package di

import (
	"reflect"
)

// asSingleton creates a singleton wrapper.
func asSingleton(provider provider) *singletonWrapper {
	return &singletonWrapper{provider: provider}
}

// singletonWrapper is a embedParamProvider wrapper. Stores provided value for prevent reinitialization.
type singletonWrapper struct {
	provider               // source provider
	value    reflect.Value // value cache
}

// Provide
func (s *singletonWrapper) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	if s.value.IsValid() {
		return s.value, nil, nil
	}
	value, cleanup, err := s.provider.Provide(values...)
	s.value = value
	return value, cleanup, err
}
