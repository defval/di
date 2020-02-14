package di

import "fmt"

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
