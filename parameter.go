package di

import (
	"fmt"
	"reflect"
)

// parameterRequired
type parameter struct {
	uniq     string       // internal uniq key
	typ      reflect.Type // resultant type
	name     string       // string identifier
	optional bool         // optional flag
}

// ID returns parameter identity.
func (p parameter) ID() key {
	return key{p.typ, p.name}
}

// String represents parameter as string.
func (p parameter) String() string {
	return p.ID().String()
}

// ResolveProvider resolves type in container c.
func (p parameter) ResolveProvider(c *Container) (provider, error) {
	plist, exists := c.providers[p.typ]
	// only one provider of type
	if exists && plist.Len() == 1 {
		return plist.ByIndex(0), nil
	}
	if exists && p.uniq != "" {
		return plist.ByUniq(p.uniq), nil
	}
	// named provider
	if exists && plist.Len() > 1 && p.name != "" {
		prov, ok := findNamedProvider(plist, p.name)
		if !ok {
			return nil, errParameterProviderNotFound{p}
		}
		return prov, nil
	}
	if exists && plist.Len() > 1 && p.name == "" {
		return nil, errHaveSeveralInstances{p.typ}
	}
	// injectable parameter
	if !exists && isInjectable(p.typ) {
		// constructor result with di.Inject - only addressable pointers
		// anonymous parameters with di.Inject - only struct
		if p.typ.Kind() == reflect.Ptr {
			return nil, errParameterProviderNotFound{p}
		}
		return providerFromInjectableParameter(p), nil
	}
	// not group of type
	if !exists && p.typ.Kind() != reflect.Slice {
		return nil, errParameterProviderNotFound{p}
	}
	// check group
	if !exists && p.typ.Kind() == reflect.Slice {
		gtype := p.typ.Elem()
		all, ok := c.providers[gtype]
		if !ok {
			return nil, errParameterProviderNotFound{p}
		}
		return newProviderGroup(p.typ, all), nil
	}
	return nil, errParameterProviderNotFound{p}
}

// ResolveValue resolves value in container c.
func (p parameter) ResolveValue(c *Container) (reflect.Value, error) {
	provider, err := p.ResolveProvider(c)
	if _, ok := err.(errParameterProviderNotFound); ok && p.optional {
		return reflect.New(p.typ).Elem(), nil
	}
	if err != nil {
		return reflect.Value{}, err
	}
	plist := provider.ParameterList()
	values, err := plist.Resolve(c)
	if err != nil {
		switch cerr := err.(type) {
		case errParameterProviderNotFound:
			return reflect.Value{}, errDependencyNotFound{p.ID(), cerr.param.ID()}
		default:
			return reflect.Value{}, fmt.Errorf("%s: %s", p, err)
		}
	}
	if len(plist) > 0 {
		c.logger.Logf("%s resolved with: %s", p, plist)
	} else {
		c.logger.Logf("%s resolved", p)
	}
	value, cleanup, err := provider.Provide(values...)
	if err != nil {
		return value, errParameterProvideFailed{p, err}
	}
	if cleanup != nil {
		c.cleanups = append(c.cleanups, cleanup)
	}
	return value, nil
}
