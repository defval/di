package di

import (
	"errors"
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
	name     string
	ctor     reflection.Func
	ctorType ctorType
	clean    *reflection.Func
}

// newProviderConstructor creates new constructor provider with name as additional identity key.
func newProviderConstructor(name string, constructor interface{}) (*providerConstructor, error) {
	if constructor == nil {
		return nil, fmt.Errorf("constructor must be a function like func([dep1, dep2, ...]) (<result>, [cleanup, error]), got nil")
	}
	fn, isFunc := reflection.InspectFunc(constructor)
	if !isFunc {
		return nil, fmt.Errorf("constructor must be a function like func([dep1, dep2, ...]) (<result>, [cleanup, error]), got %s", reflect.TypeOf(constructor))
	}
	ctorType := determineCtorType(fn)
	if ctorType == ctorUnknown {
		return nil, fmt.Errorf("constructor must be a function like func([dep1, dep2, ...]) (<result>, [cleanup, error]), got %s", reflect.TypeOf(constructor))
	}
	return &providerConstructor{
		name:     name,
		ctor:     fn,
		ctorType: determineCtorType(fn),
	}, nil
}

// ID returns provider resultant type id.
func (c providerConstructor) ID() id {
	return id{
		Name: c.name,
		Type: c.ctor.Out(0),
	}
}

// ParameterList returns constructor parameter list.
func (c providerConstructor) ParameterList() parameterList {
	var plist parameterList
	for i := 0; i < c.ctor.NumIn(); i++ {
		typ := c.ctor.In(i)
		p := parameter{
			// name:     "",
			typ:      typ,
			optional: false,
			embed:    isEmbedParameter(typ),
		}
		plist = append(plist, p)
	}
	return plist
}

// Provide provides resultant.
func (c *providerConstructor) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	out := reflection.CallResult(c.ctor.Call(values))
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
	return reflect.Value{}, nil, errors.New("you found a bug, please create new issue for " +
		"this: https://github.com/goava/di/issues/new")
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
