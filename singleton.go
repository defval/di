package di

import "reflect"

type providerSingleton struct {
	provider
	instance reflect.Value
}

func asSingleton(p provider) *providerSingleton {
	return &providerSingleton{provider: p}
}

func (p *providerSingleton) Provide(values ...reflect.Value) (_ reflect.Value, cleanup func(), err error) {
	if p.instance.IsValid() {
		return p.instance, nil, nil
	}
	value, cleanup, err := p.provider.Provide(values...)
	if err != nil {
		return value, cleanup, err
	}
	p.instance = value
	return value, cleanup, err
}
