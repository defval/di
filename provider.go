package di

import (
	"fmt"
	"reflect"
	"strings"
)

// provider is a internal provider interface.
type provider interface {
	// Type returns result type of provider.
	Type() reflect.Type
	// Name return optional string identifier.
	Name() string
	// ParameterList returns array of dependencies.
	ParameterList() parameterList
	// Provide provides value from provided parameters.
	Provide(values ...reflect.Value) (_ reflect.Value, cleanup func(), err error)
}

type uniqueProvider struct {
	provider
	uniq string
}

func createProviderList() *providerList {
	return &providerList{
		order:     []string{},
		providers: map[string]provider{},
	}
}

type providerList struct {
	order     []string
	providers map[string]provider
}

// Add adds provider to list and return uniq identifier. If named provider already
// exists in list returns error.
func (l *providerList) Add(p provider, tags map[string]string) (string, error) {
	for _, prov := range l.providers {
		if p.Name() != "" && p.Name() == prov.Name() {
			return "", fmt.Errorf("%s with name %s already exists, use another name", p.Type(), p.Name())
		}
	}
	uniq := randomString(32)
	if tags != nil {
		// identity=foo:value;bar:value
		uniq = uniq + "=" + tagsToString(tags)
	}
	l.order = append(l.order, uniq)
	l.providers[uniq] = p
	return uniq, nil
}

// All return all providers as is.
func (l providerList) All() (result []provider) {
	for _, uniq := range l.order {
		result = append(result, l.providers[uniq])
	}
	return result
}

// Uniqs returns all providers with uniq.
func (l providerList) Uniqs() (result []uniqueProvider) {
	for _, uniq := range l.order {
		result = append(result, uniqueProvider{
			uniq:     uniq,
			provider: l.providers[uniq],
		})
	}
	return result
}

func (l providerList) Len() int {
	return len(l.order)
}

func (l providerList) ByIndex(index int) provider {
	return l.providers[l.order[index]]
}

func (l providerList) ByUniq(uniq string) (provider, error) {
	for providerUniq := range l.providers {
		if strings.Contains(providerUniq, uniq) {
			return l.providers[providerUniq], nil
		}
	}
	return nil, errTaggedTypeNotFound{uniq}
}

func findNamedProvider(plist *providerList, param parameter) (result provider, _ error) {
	for _, p := range plist.All() {
		if p.Name() == param.name {
			if result != nil {
				return nil, errHaveSeveralInstances{typ: param.typ}
			}
			result = p
		}
	}
	if param.name == "" && result == nil {
		return nil, errHaveSeveralInstances{typ: param.typ}
	}
	if result != nil {
		return result, nil
	}
	return nil, errParameterProviderNotFound{param: param}
}
