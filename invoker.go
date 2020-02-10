package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/reflection"
)

type invokerType int

const (
	invokerUnknown invokerType = iota
	invokerStd                 // func (deps) {}
	invokerError               // func (deps) error {}
)

func determineInvokerType(fn *reflection.Func) (invokerType, error) {
	if fn.NumOut() == 0 {
		return invokerStd, nil
	}
	if fn.NumOut() == 1 && reflection.IsError(fn.Out(0)) {
		return invokerError, nil
	}
	return invokerUnknown, fmt.Errorf("the addInvocation function must be a function like `func([dep1, dep2, ...]) [error]`, got `%s`", fn.Type)
}

type invoker struct {
	typ invokerType
	fn  *reflection.Func
}

func newInvoker(fn interface{}) (*invoker, error) {
	if fn == nil {
		return nil, fmt.Errorf("the addInvocation function must be a function like `func([dep1, dep2, ...]) [error]`, got `%s`", "nil")
	}
	if !reflection.IsFunc(fn) {
		return nil, fmt.Errorf("the addInvocation function must be a function like `func([dep1, dep2, ...]) [error]`, got `%s`", reflect.ValueOf(fn).Type())
	}
	ifn := reflection.InspectFunction(fn)
	typ, err := determineInvokerType(ifn)
	if err != nil {
		return nil, err
	}
	return &invoker{
		typ: typ,
		fn:  reflection.InspectFunction(fn),
	}, nil
}

func (i *invoker) Invoke(c *Container) error {
	plist := i.parameters()
	values, err := plist.Resolve(c)
	if err != nil {
		return fmt.Errorf("resolve invocation (%s): %s", i.fn.Name, err)
	}
	results := i.fn.Call(values)
	if len(results) == 0 {
		return nil
	}
	if results[0].Interface() == nil {
		return nil
	}
	return results[0].Interface().(error)
}

func (i *invoker) parameters() parameterList {
	var plist parameterList
	for j := 0; j < i.fn.NumIn(); j++ {
		ptype := i.fn.In(j)
		p := parameter{
			res:      ptype,
			optional: false,
			embed:    isEmbedParameter(ptype),
		}
		plist = append(plist, p)
	}
	return plist
}
