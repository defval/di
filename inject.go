package di

import (
	"reflect"
	"strconv"
)

// Inject indicates that struct public fields will be injected automatically.
//
//	type Application struct {
//		di.Inject
//
//		Server *http.Server // will be injected
//	}
//
// You can specify tags for injected types:
//
//  type Application struct {
//  	di.Inject
//
//		Public 	*http.Server `type:"public"` 	// *http.Server with type:public tag combination will be injected
//		Private *http.Server `type:"private"` 	// *http.Server with type:private tag combination will be injected
//  }
type Inject struct {
	injectable
}

// injectable interface needs to struct fields injection functional.
type injectable interface {
	isInjectable()
}

type field struct {
	rt       reflect.Type
	tags     Tags
	optional bool
}

// canInject checks that type t contain di.Inject and supports injecting.
func canInject(t reflect.Type) bool {
	if !t.Implements(injectableInterface) {
		return false
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	return true
}

// parsePopulateFields parses fields of struct that can be populated.
func parsePopulateFields(rt reflect.Type) map[int]field {
	if !canInject(rt) {
		return nil
	}
	var rv reflect.Value
	if !rv.IsValid() {
		switch rt.Kind() {
		case reflect.Ptr:
			rv = reflect.New(rt.Elem())
		default:
			rv = reflect.New(rt).Elem()
		}
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	fields := make(map[int]field, rt.NumField())
	// fi - field index
	for fi := 0; fi < rt.NumField(); fi++ {
		fv := rv.Field(fi)
		// check that field can be set
		if !fv.CanSet() {
			continue
		}
		// cur - current field
		cur := rt.Field(fi)
		f, valid := inspectStructField(cur)
		if !valid {
			continue
		}
		fields[fi] = field{
			rt:       cur.Type,
			tags:     f.tags,
			optional: f.optional,
		}
	}
	return fields
}

// inspectStructField parses struct field
func inspectStructField(f reflect.StructField) (field, bool) {
	tags := Tags{}
	t := string(f.Tag)
	optional := false
	// this code copied from reflect.StructField.Lookup() method.
	for t != "" {
		// Skip leading space.
		i := 0
		for i < len(t) && t[i] == ' ' {
			i++
		}
		t = t[i:]
		if t == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(t) && t[i] > ' ' && t[i] != ':' && t[i] != '"' && t[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(t) || t[i] != ':' || t[i+1] != '"' {
			break
		}
		name := string(t[:i])
		t = t[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(t) && t[i] != '"' {
			if t[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(t) {
			break
		}
		qvalue := string(t[:i+1])
		t = t[i+1:]
		value, err := strconv.Unquote(qvalue)
		if err != nil {
			break
		}
		if name == "skip" && value == "true" {
			return field{}, false
		}
		if name == "optional" {
			if value == "true" {
				optional = true
			}
			continue
		}
		tags[name] = value
	}
	return field{
		rt:       f.Type,
		tags:     tags,
		optional: optional,
	}, true
}

var injectableInterface = reflect.TypeOf(new(injectable)).Elem()
