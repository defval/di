package di

import (
	"fmt"

	"github.com/goava/di/internal/stacktrace"
)

// ErrParameterProvideFailed causes when container found a provider but provide failed.
type ErrParameterProvideFailed struct {
	id  id    // type identity
	err error // error
}

// Error is a implementation of error interface.
func (e ErrParameterProvideFailed) Error() string {
	return fmt.Sprintf("%s: %s", e.id, e.err)
}

// ErrParameterProviderNotFound causes when container could not found a provider for parameter.
type ErrParameterProviderNotFound struct {
	param parameter
}

// Error is a implementation of error interface.
func (e ErrParameterProviderNotFound) Error() string {
	return fmt.Sprintf("%s: not exists in container", e.param)
}

type errProvideFailed []struct {
	frame stacktrace.Frame
	err   error
}

// Append returns new errProvideFailed with appended error.
func (e errProvideFailed) Append(frame stacktrace.Frame, err error) errProvideFailed {
	return append(e, struct {
		frame stacktrace.Frame
		err   error
	}{frame: frame, err: err})
}

// String implements Stinger interface.
func (e errProvideFailed) String() string {
	var str string
	for _, entry := range e {
		str += fmt.Sprintf("\t%s: %s\n", entry.frame, entry.err)
	}
	return "di.Provide(..) failed:\n" + str
}

// Error implements error interface.
func (e errProvideFailed) Error() string {
	return e.String()
}

type errInvokeFailed struct {
	frame stacktrace.Frame
	err   error
}

// Error implements error interface.
func (e errInvokeFailed) Error() string {
	return fmt.Sprintf("di.Invoke(..) failed: %s: %s", e.frame, e.err)
}

type errResolveFailed struct {
	frame stacktrace.Frame
	err   error
}

// Error implements error interface.
func (e errResolveFailed) Error() string {
	return fmt.Sprintf("di.Resolve(..) failed: %s: %s", e.frame, e.err)
}
