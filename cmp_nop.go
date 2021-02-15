package di

import (
	"reflect"
)

type nopCompiler struct{}

func (v nopCompiler) deps(s schema) ([]*node, error) {
	return nil, nil
}

func (v nopCompiler) compile(dependencies []reflect.Value, s schema) (reflect.Value, error) {
	bug()
	return reflect.Value{}, nil
}
