package di

import (
	"reflect"
)

// parameterList
type parameterList []parameter

// ResolveValues loads all parameters presented in parameter list.
func (pl parameterList) Resolve(c *Container) ([]reflect.Value, error) {
	var values []reflect.Value
	for _, p := range pl {
		value, err := p.ResolveValue(c)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}
