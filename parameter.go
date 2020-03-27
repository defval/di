package di

import (
	"reflect"
	"strings"
	"unsafe"
)

// isEmbedParameter
func isEmbedParameter(typ reflect.Type) bool {
	return typ.Kind() == reflect.Struct && typ.Implements(parameterInterface)
}

// parameterRequired
type parameter struct {
	name     string       // string identifier
	typ      reflect.Type // resultant type
	optional bool         // optional flag
	embed    bool         // embed flag
}

// String represents parameter as string.
func (p parameter) String() string {
	return id{Name: p.name, Type: p.typ}.String()
}

// ResolveProvider resolves type in container c.
func (p parameter) ResolveProvider(c *Container) (provider, bool) {
	k := id{
		Name: p.name,
		Type: p.typ,
	}
	node, err := c.graph.Node(k)
	if err != nil {
		return nil, false
	}
	return node.(providerNode).provider, true
}

// ResolveValue resolves value in container c.
func (p parameter) ResolveValue(c *Container) (reflect.Value, error) {
	provider, exists := p.ResolveProvider(c)
	if !exists && p.optional {
		return reflect.New(p.typ).Elem(), nil
	}
	if !exists {
		return reflect.Value{}, ErrParameterProviderNotFound{param: p}
	}
	pl := provider.ParameterList()
	if len(pl) > 0 {
		c.logger.Logf("%s resolved with: %s", p, pl)
	} else {
		c.logger.Logf("%s resolved", p)
	}
	values, err := pl.Resolve(c)
	if err != nil {
		return reflect.Value{}, err
	}
	value, cleanup, err := provider.Provide(values...)
	if err != nil {
		return value, ErrParameterProvideFailed{id: provider.ID(), err: err}
	}
	if cleanup != nil {
		c.cleanups = append(c.cleanups, cleanup)
	}
	// inject struct
	err = p.ResolveProperty(c, value)
	if err != nil {
		return reflect.Value{}, err
	}
	return value, nil
}

const (
	flagRO = 0b1100000
)

func ValuePatch(v reflect.Value) reflect.Value {
	rv := reflect.ValueOf(&v)
	flag := rv.Elem().FieldByName("flag")
	ptrFlag := (*uintptr)(unsafe.Pointer(flag.UnsafeAddr()))
	*ptrFlag = *ptrFlag &^ flagRO
	return v
}

func (p parameter) ResolveProperty(c *Container, value reflect.Value) (err error) {
	value = reflect.Indirect(value)
	if value.Kind() != reflect.Struct {
		return nil
	}
	vType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		fieldType := vType.Field(i)
		field := ValuePatch(value.Field(i))
		tag, ok := fieldType.Tag.Lookup("di")
		if ok {
			var optional = false
			var name string
			tags := strings.Split(tag, ",")
			for _, t := range tags {
				t = strings.Trim(t, " ")
				if t == "optional" {
					optional = true
					continue
				}
				name = t
			}
			pp := parameter{
				name:     name,
				typ:      fieldType.Type,
				optional: optional,
			}
			param, err := pp.ResolveValue(c)
			if err != nil {
				return err
			}
			field.Set(param)
		}
	}
	return nil
}

// internalParameter
type internalParameter interface {
	isDependencyInjectionParameter()
}

// parameterInterface
var parameterInterface = reflect.TypeOf(new(internalParameter)).Elem()
