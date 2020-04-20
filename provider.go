package di

import (
	"fmt"
	"reflect"
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

func (l *providerList) Add(p provider) error {
	for _, prov := range l.providers {
		if p.Name() != "" && p.Name() == prov.Name() {
			return fmt.Errorf("%s with name %s already exists, use another name", p.Type(), p.Name())
		}
	}
	uniq := randomString(32)
	l.order = append(l.order, uniq)
	l.providers[uniq] = p
	return nil
}

func (l providerList) All() (result []provider) {
	for _, uniq := range l.order {
		result = append(result, l.providers[uniq])
	}
	return result
}

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

func (l providerList) ByUniq(uniq string) provider {
	return l.providers[uniq]
}

func findNamedProvider(plist *providerList, name string) (provider, bool) {
	for _, p := range plist.All() {
		if p.Name() == name {
			return p, true
		}
	}
	return nil, false
}
