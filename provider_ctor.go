package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/reflection"
)

// ctorType describes types of constructor provider.
type ctorType int

const (
	ctorUnknown      ctorType = iota
	ctorStd                   // (deps) (result)
	ctorError                 // (deps) (result, error)
	ctorCleanup               // (deps) (result, cleanup)
	ctorCleanupError          // (deps) (result, cleanup, error)
)

// providerConstructor is a provider that can handle constructor functions.
// Type of this provider provides type with function call.
type providerConstructor struct {
	name string
	// constructor
	ctorType ctorType
	call     reflection.Func
	// injectable params
	injectable struct {
		// params parsed once
		params []parameter
		// field numbers parsed once
		fields []int
	}
	// clean callback
	clean *reflection.Func
}

// newProviderConstructor creates new constructor provider with name as additional identity key.
func newProviderConstructor(name string, constructor interface{}) (*providerConstructor, error) {
	if constructor == nil {
		return nil, fmt.Errorf("constructor must be a function like func([dep1, dep2, ...]) (<result>, [cleanup, error]), got nil")
	}
	fn, isFn := reflection.InspectFunc(constructor)
	if !isFn {
		return nil, fmt.Errorf("constructor must be a function like func([dep1, dep2, ...]) (<result>, [cleanup, error]), got %s", reflect.TypeOf(constructor))
	}
	ctorType := determineCtorType(fn)
	if ctorType == ctorUnknown {
		return nil, fmt.Errorf("constructor must be a function like func([dep1, dep2, ...]) (<result>, [cleanup, error]), got %s", reflect.TypeOf(constructor))
	}
	provider := &providerConstructor{
		name:     name,
		call:     fn,
		ctorType: ctorType,
	}
	// result type
	rt := fn.Out(0)
	// if struct is injectable, range over injectableFields and parse injectable params
	if isInjectable(rt) {
		provider.injectable.params, provider.injectable.fields = parseInjectableType(rt)
	}
	return provider, nil
}

// ID returns provider resultant type id.
func (c providerConstructor) ID() id {
	return id{
		Name: c.name,
		Type: c.call.Out(0),
	}
}

// ParameterList returns constructor parameter list.
func (c *providerConstructor) ParameterList() parameterList {
	// todo: move to constructor
	var pl parameterList
	for i := 0; i < c.call.NumIn(); i++ {
		in := c.call.In(i)
		pl = append(pl, parameter{
			name: "", // constructor parameters could be resolved only with empty ("") name
			typ:  in,
		})
	}
	return append(pl, c.injectable.params...)
}

// Provide provides resultant.
func (c *providerConstructor) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	// constructor last param index
	clpi := c.call.NumIn()
	if c.call.NumIn() == 0 {
		clpi = 0
	}
	out := reflection.CallResult(c.call.Call(values[:clpi]))
	if c.ctorType == ctorError && out.Error(1) != nil {
		return out.Result(), nil, out.Error(1)
	}
	if c.ctorType == ctorCleanupError && out.Error(2) != nil {
		return out.Result(), nil, out.Error(2)
	}
	// set injectable fields
	if len(c.injectable.fields) > 0 {
		// result value
		rv := out.Result()
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		fields := values[clpi:]
		// field index
		for i, value := range fields {
			// field value
			fv := rv.Field(c.injectable.fields[i])
			if !fv.CanSet() {
				panic("you found a bug, please create new issue for this: https://github.com/goava/di/issues/new")
			}
			fv.Set(value)
		}
	}
	switch c.ctorType {
	case ctorStd:
		return out.Result(), nil, nil
	case ctorError:
		return out.Result(), nil, out.Error(1)
	case ctorCleanup:
		return out.Result(), out.Cleanup(), nil
	case ctorCleanupError:
		return out.Result(), out.Cleanup(), out.Error(2)
	}
	panic("you found a bug, please create new issue for this: https://github.com/goava/di/issues/new")
}

// determineCtorType
func determineCtorType(fn reflection.Func) ctorType {
	if fn.NumOut() == 1 {
		return ctorStd
	}
	if fn.NumOut() == 2 {
		if reflection.IsError(fn.Out(1)) {
			return ctorError
		}
		if reflection.IsCleanup(fn.Out(1)) {
			return ctorCleanup
		}
	}
	if fn.NumOut() == 3 && reflection.IsCleanup(fn.Out(1)) && reflection.IsError(fn.Out(2)) {
		return ctorCleanupError
	}
	return ctorUnknown
}
