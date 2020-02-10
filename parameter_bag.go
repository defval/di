package di

import (
	"fmt"
	"reflect"
)

// createParameterBugProvider
func createParameterBugProvider(key key, parameters ParameterBag) internalProvider {
	return newProviderConstructor(key.String(), func() ParameterBag { return parameters })
}

// parameterBagType
var parameterBagType = reflect.TypeOf(ParameterBag{})

// ParameterBag
type ParameterBag map[string]interface{}

// Exists
func (b ParameterBag) Exists(key string) bool {
	_, ok := b[key]
	return ok
}

// Get
func (b ParameterBag) Get(key string) (interface{}, bool) {
	value, ok := b[key]
	return value, ok
}

// String
func (b ParameterBag) String(key string) (string, bool) {
	value, ok := b[key].(string)
	return value, ok
}

// Int64
func (b ParameterBag) Int64(key string) (int64, bool) {
	value, ok := b[key].(int64)
	return value, ok
}

// Int
func (b ParameterBag) Int(key string) (int, bool) {
	value, ok := b[key].(int)
	return value, ok
}

// Float64
func (b ParameterBag) Float64(key string) (float64, bool) {
	value, ok := b[key].(float64)
	return value, ok
}

// Require
func (b ParameterBag) Require(key string) interface{} {
	value, ok := b[key]
	if !ok {
		panic(fmt.Sprintf("value for string key `%s` not found", key))
	}
	return value
}

// RequireString
func (b ParameterBag) RequireString(key string) string {
	value, ok := b[key].(string)
	if !ok {
		panic(fmt.Sprintf("value for string key `%s` not found", key))
	}
	return value
}

// RequireInt64
func (b ParameterBag) RequireInt64(key string) int64 {
	value, ok := b[key].(int64)
	if !ok {
		panic(fmt.Sprintf("value for string key `%s` not found", key))
	}
	return value
}

// RequireInt
func (b ParameterBag) RequireInt(key string) int {
	value, ok := b[key].(int)
	if !ok {
		panic(fmt.Sprintf("value for string key `%s` not found", key))
	}
	return value
}

// RequireFloat64
func (b ParameterBag) RequireFloat64(key string) float64 {
	value, ok := b[key].(float64)
	if !ok {
		panic(fmt.Sprintf("value for string key `%s` not found", key))
	}
	return value
}
