package di

import (
	"reflect"
)

// parameterRequired
type parameter struct {
	name     string       // string identifier
	typ      reflect.Type // resultant type
	optional bool         // optional flag
}

// String represents parameter as string.
func (p parameter) String() string {
	return id{Name: p.name, Type: p.typ}.String()
}

// ResolveProvider resolves type in container c.
func (p parameter) ResolveProvider(c *Container) (provider, bool) {
	k := id{
		Name: p.name,
		Type: p.typ,
	}
	node, err := c.graph.Node(k)
	if err != nil {
		return nil, false
	}
	return node.(providerNode).provider, true
}

// ResolveValue resolves value in container c.
func (p parameter) ResolveValue(c *Container) (reflect.Value, error) {
	provider, exists := p.ResolveProvider(c)
	if !exists && p.optional {
		return reflect.New(p.typ).Elem(), nil
	}
	if !exists {
		return reflect.Value{}, ErrParameterProviderNotFound{param: p}
	}
	pl := provider.ParameterList()
	if len(pl) > 0 {
		c.logger.Logf("%s resolved with: %s", p, pl)
	} else {
		c.logger.Logf("%s resolved", p)
	}
	values, err := pl.Resolve(c)
	if err != nil {
		return reflect.Value{}, err
	}
	value, cleanup, err := provider.Provide(values...)
	if err != nil {
		return value, ErrParameterProvideFailed{id: provider.ID(), err: err}
	}
	if cleanup != nil {
		c.cleanups = append(c.cleanups, cleanup)
	}
	return value, nil
}
