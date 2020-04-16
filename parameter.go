package di

import (
	"fmt"
	"reflect"
)

// parameterRequired
type parameter struct {
	name     string       // string identifier
	typ      reflect.Type // resultant type
	optional bool         // optional flag
}

// ID returns parameter identity.
func (p parameter) ID() id {
	return id{
		Name: p.name,
		Type: p.typ,
	}
}

// String represents parameter as string.
func (p parameter) String() string {
	return id{Name: p.name, Type: p.typ}.String()
}

// ResolveProvider resolves type in container c.
func (p parameter) ResolveProvider(c *Container) (provider, bool) {
	id := id{
		Name: p.name,
		Type: p.typ,
	}
	provider, exists := c.providers[id]
	if !exists {
		return nil, false
	}
	return provider, true
}

// ResolveValue resolves value in container c.
func (p parameter) ResolveValue(c *Container) (reflect.Value, error) {
	_, prototype := c.prototypes[p.ID()]
	if existing, ok := c.values[p.ID()]; ok && !prototype {
		return existing, nil
	}
	provider, exists := p.ResolveProvider(c)
	if !exists && isInjectable(p.typ) {
		exists = true
		provider = providerFromInjectableParameter(p)
	}
	if !exists && p.optional {
		return reflect.New(p.typ).Elem(), nil
	}
	if !exists {
		return reflect.Value{}, ErrParameterProviderNotFound{param: p}
	}
	pl := provider.ParameterList()
	values, err := pl.Resolve(c)
	if err != nil {
		switch cerr := err.(type) {
		case ErrParameterProviderNotFound:
			return reflect.Value{}, fmt.Errorf("%s: dependency %s not exists in container", p, cerr.param)
		default:
			return reflect.Value{}, fmt.Errorf("%s: %s", p, err)
		}
	}
	if len(pl) > 0 {
		c.logger.Logf("%s resolved with: %s", p, pl)
	} else {
		c.logger.Logf("%s resolved", p)
	}
	value, cleanup, err := provider.Provide(values...)
	if err != nil {
		return value, ErrParameterProvideFailed{id: provider.ID(), err: err}
	}
	c.values[provider.ID()] = value
	if cleanup != nil {
		c.cleanups = append(c.cleanups, cleanup)
	}
	return value, nil
}
