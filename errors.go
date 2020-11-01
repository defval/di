package di

import (
	"fmt"

	"github.com/goava/di/internal/stacktrace"
)

// knownError return true if err is library known error.
func knownError(err error) bool {
	switch err.(type) {
	case errProvideFailed, errResolveFailed, errInvokeFailed, errInvalidInvocation:
		return true
	default:
		return false
	}
}

func provideErrWithStack(err error) errProvideFailed {
	return errProvideFailed{stacktrace.CallerFrame(1), err}
}

func invokeErrWithStack(err error) errInvokeFailed {
	return errInvokeFailed{stacktrace.CallerFrame(1), err}
}

// errProvideFailed causes when constructor providing failed.
type errProvideFailed struct {
	frame stacktrace.Frame
	err   error
}

// Error returns error string.
func (e errProvideFailed) Error() string {
	return fmt.Sprintf("%s: %s", e.frame, e.err)
}

// errInvokeFailed causes when invoke failed.
type errInvokeFailed struct {
	frame stacktrace.Frame
	err   error
}

// Error returns error string.
func (e errInvokeFailed) Error() string {
	return fmt.Sprintf("%s: %s", e.frame, e.err)
}

func resolveErrWithStack(err error) errResolveFailed {
	return errResolveFailed{
		frame: stacktrace.CallerFrame(1),
		err:   err,
	}
}

// errResolveFailed causes when type resolve failed.
type errResolveFailed struct {
	frame stacktrace.Frame
	err   error
}

// Error returns error string.
func (e errResolveFailed) Error() string {
	return fmt.Sprintf("%s: %s", e.frame, e.err)
}

type errInvalidInvocation struct {
	err error
}

func (e errInvalidInvocation) Error() string {
	return e.err.Error()
}

func bug() {
	panic("you found a bug, please create new issue for this: https://github.com/goava/di/issues/new")
}
