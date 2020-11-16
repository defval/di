package di

import (
	"errors"
	"fmt"

	"github.com/goava/di/internal/stacktrace"
)

var (
	// ErrTypeNotExists causes when type not found in container.
	ErrTypeNotExists = errors.New("not exists in the container")
)

var (
	errInvalidInvocationSignature = errors.New("invalid invocation signature")
	errCycleDetected              = errors.New("cycle detected")
	errFieldsNotSupported         = errors.New("fields not supported")
)

// knownError return true if err is library known error.
func knownError(err error) bool {
	if errors.Is(err, ErrTypeNotExists) ||
		errors.Is(err, errInvalidInvocationSignature) ||
		errors.Is(err, errCycleDetected) ||
		errors.Is(err, errFieldsNotSupported) {
		return true
	}
	return false
}

func errWithStack(err error) error {
	return fmt.Errorf("%s: %w", stacktrace.CallerFrame(1), err)
}

func bug() {
	panic("you found a bug, please create new issue for this: https://github.com/goava/di/issues/new")
}
