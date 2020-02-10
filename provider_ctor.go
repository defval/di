package di

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/goava/di/internal/reflection"
)

type ctorType int

const (
	ctorUnknown      ctorType = iota // unknown ctor signature
	ctorStd                          // (deps) (result)
	ctorError                        // (deps) (result, error)
	ctorCleanup                      // (deps) (result, cleanup)
	ctorCleanupError                 // (deps) (result, cleanup, error)
)

// newProviderConstructor
func newProviderConstructor(name string, ctor interface{}) *providerConstructor {
	if ctor == nil {
		panicf("The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `%s`", "nil")
	}
	if !reflection.IsFunc(ctor) {
		panicf("The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `%s`", reflect.ValueOf(ctor).Type())
	}
	fn := reflection.InspectFunction(ctor)
	ctorType := determineCtorType(fn)
	return &providerConstructor{
		name:     name,
		ctor:     fn,
		ctorType: ctorType,
	}
}

// providerConstructor
type providerConstructor struct {
	name     string
	ctor     *reflection.Func
	ctorType ctorType
	clean    *reflection.Func
}

func (c providerConstructor) Key() key {
	return key{
		name: c.name,
		res:  c.ctor.Out(0),
		typ:  ptConstructor,
	}
}

func (c providerConstructor) ParameterList() parameterList {
	var plist parameterList
	for i := 0; i < c.ctor.NumIn(); i++ {
		ptype := c.ctor.In(i)
		var name string
		if ptype == parameterBagType {
			name = c.Key().String()
		}
		p := parameter{
			name:     name,
			res:      ptype,
			optional: false,
			embed:    isEmbedParameter(ptype),
		}
		plist = append(plist, p)
	}
	return plist
}

// Provide
func (c *providerConstructor) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	out := callResult(c.ctor.Call(values))
	switch c.ctorType {
	case ctorStd:
		return out.instance(), nil, nil
	case ctorError:
		return out.instance(), nil, out.error(1)
	case ctorCleanup:
		return out.instance(), out.cleanup(), nil
	case ctorCleanupError:
		return out.instance(), out.cleanup(), out.error(2)
	}
	return reflect.Value{}, nil, errors.New("you found a bug, please create new issue for " +
		"this: https://github.com/goava/di/issues/new")
}

// determineCtorType
func determineCtorType(fn *reflection.Func) ctorType {
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
	panic(fmt.Sprintf("The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `%s`", fn.Name))
}

// callResult
type callResult []reflect.Value

func (r callResult) instance() reflect.Value {
	return r[0]
}

func (r callResult) cleanup() func() {
	if r[1].IsNil() {
		return nil
	}
	return r[1].Interface().(func())
}

func (r callResult) error(position int) error {
	if r[position].IsNil() {
		return nil
	}
	return r[position].Interface().(error)
}
