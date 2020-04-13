package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/reflection"
)

type invocationType int

const (
	invocationUnknown invocationType = iota
	invokerStd                       // func (deps) {}
	invokerError                     // func (deps) error {}
)

func determineInvokerType(fn reflection.Func) (invocationType, error) {
	if fn.NumOut() == 0 {
		return invokerStd, nil
	}
	if fn.NumOut() == 1 && reflection.IsError(fn.Out(0)) {
		return invokerError, nil
	}
	return invocationUnknown, fmt.Errorf("invoke function must be a function like `func([dep1, dep2, ...]) [error]`, got %s", fn.Type)
}

type invoker struct {
	typ invocationType
	fn  reflection.Func
}

func newInvoker(fn interface{}) (*invoker, error) {
	if fn == nil {
		return nil, fmt.Errorf("invoke function must be a function like `func([dep1, dep2, ...]) [error]`, got %s", "nil")
	}
	inspected, isFn := reflection.InspectFunc(fn)
	if !isFn {
		return nil, fmt.Errorf("invoke function must be a function like `func([dep1, dep2, ...]) [error]`, got %s", reflect.TypeOf(fn))
	}
	typ, err := determineInvokerType(inspected)
	if err != nil {
		return nil, err
	}
	return &invoker{
		typ: typ,
		fn:  inspected,
	}, nil
}

func (i *invoker) Invoke(c *Container) error {
	plist := i.parameters()
	values, err := plist.Resolve(c)
	if err != nil {
		return fmt.Errorf("resolve invocation (%s): %s", i.fn.Name, err)
	}
	results := reflection.CallResult(i.fn.Call(values))
	if len(results) == 0 {
		return nil
	}
	return results.Error(0)
}

func (i *invoker) parameters() parameterList {
	var plist parameterList
	for j := 0; j < i.fn.NumIn(); j++ {
		ptype := i.fn.In(j)
		p := parameter{
			typ:      ptype,
			optional: false,
		}
		plist = append(plist, p)
	}
	return plist
}
