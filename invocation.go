package di

import (
	"github.com/goava/di/internal/reflection"
)

func validateInvocation(fn reflection.Func) bool {
	if fn.NumOut() == 0 {
		return true
	}
	if fn.NumOut() == 1 && reflection.IsError(fn.Out(0)) {
		return true
	}
	return false
}

type errInvalidInvocation struct {
	err error
}

func (e errInvalidInvocation) Error() string {
	return e.err.Error()
}
