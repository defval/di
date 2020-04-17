package di

import (
	"fmt"

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

func newInvoker(fn reflection.Func) (*invoker, error) {
	typ, err := determineInvokerType(fn)
	if err != nil {
		return nil, err
	}
	return &invoker{
		typ: typ,
		fn:  fn,
	}, nil
}

func (i *invoker) Invoke(c *Container) error {
	plist := i.parameters()
	values, err := plist.Resolve(c)
	if err != nil {
		return err
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
