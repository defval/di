package graphkv

import "fmt"

// ErrKeyAlreadyExists
type ErrKeyAlreadyExists struct {
	Key Key
}

func (e ErrKeyAlreadyExists) Error() string {
	return fmt.Sprintf("%s already exists", e.Key)
}

// ErrNodeNotExists
type ErrNodeNotExists struct {
	Key Key
}

// ErrNodeNotExists
func (e ErrNodeNotExists) Error() string {
	return fmt.Sprintf("%s not exists", e.Key)
}
