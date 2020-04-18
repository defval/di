package di

import (
	"fmt"

	"github.com/goava/di/internal/stacktrace"
)

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
	id  id    // type identity
	err error // error
}

// Error is a implementation of error interface.
func (e errParameterProvideFailed) Error() string {
	return fmt.Sprintf("%s: %s", e.id, e.err)
}

// errParameterProviderNotFound causes when container could not found a provider for parameter.
type errParameterProviderNotFound struct {
	param parameter
}

// Error is a implementation of error interface.
func (e errParameterProviderNotFound) Error() string {
	return fmt.Sprintf("type %s not exists in container", e.param)
}
