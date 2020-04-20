package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/stacktrace"
)

// isUsageError return true if err is library usage error.
func isUsageError(err error) bool {
	switch err.(type) {
	case
		ErrProvideFailed,
		ErrResolveFailed,
		ErrInvokeFailed,
		errInvalidInvocation,
		errParameterProviderNotFound,
		errParameterProvideFailed,
		errHaveSeveralInstances:
		return true
	default:
		return false
	}
}

func provideErrWithStack(err error) ErrProvideFailed {
	return ErrProvideFailed{stacktrace.CallerFrame(1), err}
}

// ErrProvideFailed causes when constructor providing failed.
type ErrProvideFailed struct {
	frame stacktrace.Frame
	err   error
}

// Error returns error string.
func (e ErrProvideFailed) Error() string {
	return fmt.Sprintf("%s: %s", e.frame, e.err)
}

func invokeErrWithStack(err error) ErrInvokeFailed {
	return ErrInvokeFailed{stacktrace.CallerFrame(1), err}
}

// ErrInvokeFailed causes when invoke failed.
type ErrInvokeFailed struct {
	frame stacktrace.Frame
	err   error
}

// Error returns error string.
func (e ErrInvokeFailed) Error() string {
	return fmt.Sprintf("%s: %s", e.frame, e.err)
}

func resolveErrWithStack(err error) ErrResolveFailed {
	return ErrResolveFailed{
		frame: stacktrace.CallerFrame(1),
		err:   err,
	}
}

// ErrResolveFailed causes when type resolve failed.
type ErrResolveFailed struct {
	frame stacktrace.Frame
	err   error
}

// Error returns error string.
func (e ErrResolveFailed) Error() string {
	return fmt.Sprintf("%s: %s", e.frame, e.err)
}

// errParameterProvideFailed causes when container found a provider but provide failed.
type errParameterProvideFailed struct {
	parameter parameter
	err       error // error
}

// Error is a implementation of error interface.
func (e errParameterProvideFailed) Error() string {
	return fmt.Sprintf("%s: %s", e.parameter, e.err)
}

// errParameterProviderNotFound causes when container could not found a provider for parameter.
type errParameterProviderNotFound struct {
	param parameter
}

// Error is a implementation of error interface.
func (e errParameterProviderNotFound) Error() string {
	return fmt.Sprintf("type %s not exists in container", e.param)
}

type errDependencyNotFound struct {
	dependant id
	parameter id
}

func (e errDependencyNotFound) Error() string {
	return fmt.Sprintf("%s: dependency %s not exists in container", e.dependant, e.parameter)
}

type errHaveSeveralInstances struct {
	typ reflect.Type
}

func (e errHaveSeveralInstances) Error() string {
	return fmt.Sprintf("%s: could not be resolved: have several instances", e.typ)
}

func bug() {
	panic("you found a bug, please create new issue for this: https://github.com/goava/di/issues/new")
}
