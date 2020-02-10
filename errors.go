package di

import "fmt"

// ErrParameterProvideFailed
type ErrParameterProvideFailed struct {
	k   key
	err error
}

func (e ErrParameterProvideFailed) Error() string {
	return fmt.Sprintf("%s: %s", e.k, e.err)
}

// ErrParameterProviderNotFound
type ErrParameterProviderNotFound struct {
	param parameter
}

func (e ErrParameterProviderNotFound) Error() string {
	return fmt.Sprintf("%s: not exists in container", e.param)
}
