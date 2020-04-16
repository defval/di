package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/reflection"
	"github.com/goava/di/internal/stacktrace"
)

// New creates new container with provided options. Example usage:
//
// 	func NewHTTPServer(mux *http.ServeMux) *http.Server {
// 		return &http.Server{
// 			Handler: mux,
// 		}
// 	}
// 	func NewHTTPServeMux() *http.ServeMux {
// 		return http.ServeMux{}
// 	}
//
// Container initialization code:
//
// 	container, err := di.New(
// 		di.Provide(NewHTTPServer),
// 		di.Provide(NewHTTPServeMux),
// 	)
// 	if err != nil {
//		// handle error
//	}
//	var server *http.Server
//	if err := c.Resolve(&server); err != nil {
//		// handle error
//	}
func New(options ...Option) (_ *Container, err error) {
	c := &Container{
		logger:    nopLogger{},
		providers: map[id]provider{},
		values:    map[id]reflect.Value{},
		cleanups:  []func(){},
	}
	// apply container options
	for _, opt := range options {
		opt.apply(c)
	}
	// process constructors
	provideErr := errProvideFailed{}
	for _, provide := range c.initial.provides {
		if err := c.Provide(provide.constructor, provide.options...); err != nil {
			provideErr = provideErr.Append(provide.frame, err)
		}
	}
	if len(provideErr) > 0 {
		return nil, provideErr
	}
	// error omitted because if logger could not be resolve it will be default
	_ = c.Resolve(&c.logger)
	// call initial invokes
	for _, invoke := range c.initial.invokes {
		if err := c.Invoke(invoke.fn, invoke.options...); err != nil {
			return nil, errInvokeFailed{invoke.frame, err}
		}
	}
	// initial resolves
	for _, resolve := range c.initial.resolves {
		if err := c.Resolve(resolve.target, resolve.options...); err != nil {
			return nil, errResolveFailed{resolve.frame, err}
		}
	}
	return c, nil
}

// Container is a dependency injection container.
type Container struct {
	// internal logger
	logger Logger
	// initial options
	initial struct {
		provides []provideOptions
		invokes  []invokeOptions
		resolves []resolveOptions
	}
	providers map[id]provider
	values    map[id]reflect.Value
	// cleanups
	cleanups []func()
}

// Provide provides to container reliable way to build type. The constructor will be invoked lazily on-demand.
// For more information about constructors see Constructor interface. ProvideOption can add additional behavior to
// the process of type resolving.
func (c *Container) Provide(constructor Constructor, options ...ProvideOption) (err error) {
	params := ProvideParams{}
	// apply provide options
	for _, opt := range options {
		opt.apply(&params)
	}
	// create constructor provider
	var prov provider
	if prov, err = newProviderConstructor(params.Name, constructor); err != nil {
		return err
	}
	if _, ok := c.providers[prov.ID()]; ok {
		// duplicate types not allowed
		return fmt.Errorf("%s already exists in dependency graph", prov.ID())
	}
	// add provider to graph
	c.providers[prov.ID()] = prov
	if err := c.processInterfaces(prov, params.Interfaces); err != nil {
		return err
	}
	return nil
}

func (c *Container) processInterfaces(prov provider, interfaces []Interface) error {
	iprovs := make([]*providerInterface, 0, len(interfaces))
	// process interfaces
	for _, i := range interfaces {
		// create interface provider
		iprov, err := newProviderInterface(prov, i)
		if err != nil {
			return err
		}
		_, ok := c.providers[iprov.ID()]
		if !ok {
			c.providers[iprov.ID()] = iprov
		}
		if ok {
			// todo: do not change if already stub
			// if graph node already exists, resolve it as interface restricted, but it can exists in group
			stub := newProviderStub(iprov.ID(), "have several implementations")
			c.providers[iprov.ID()] = stub
		}
		iprovs = append(iprovs, iprov)
	}
	// process group for interfaces
	for _, iprov := range iprovs {
		groupID := id{
			Type: reflect.SliceOf(iprov.ID().Type),
		}
		existing, ok := c.providers[groupID]
		if ok {
			// if group node already exists use it
			existing.(*providerGroup).Add(prov.ID())
		}
		if !ok {
			group := newProviderGroup(iprov.ID())
			group.Add(prov.ID())
			c.providers[groupID] = group
		}
	}
	return nil
}

// Resolve builds instance of target type and fills target pointer.
func (c *Container) Resolve(into interface{}, options ...ResolveOption) error {
	if into == nil {
		return fmt.Errorf("resolve target must be a pointer, got nil")
	}
	if !reflection.IsPtr(into) {
		return fmt.Errorf("resolve target must be a pointer, got %s", reflect.TypeOf(into))
	}
	params := ResolveParams{}
	// apply extract options
	for _, opt := range options {
		opt.apply(&params)
	}
	typ := reflect.TypeOf(into).Elem()
	param := parameter{
		name: params.Name,
		typ:  typ,
	}
	value, err := param.ResolveValue(c)
	if err != nil {
		return err
	}
	targetValue := reflect.ValueOf(into).Elem()
	targetValue.Set(value)
	return nil
}

// Invoke calls provided function.
func (c *Container) Invoke(fn Invocation, options ...InvokeOption) error {
	// params := InvokeParams{}
	// for _, opt := range options {
	// 	opt.apply(&params)
	// }
	invoker, err := newInvoker(fn)
	if err != nil {
		return err
	}
	return invoker.Invoke(c)
}

// Has checks that type exists in container, if not it return false.
func (c *Container) Has(target interface{}, options ...ResolveOption) bool {
	if target == nil {
		return false
	}
	params := ResolveParams{}
	// apply options
	for _, opt := range options {
		opt.apply(&params)
	}
	typ := reflect.TypeOf(target)
	param := parameter{
		name: params.Name,
		typ:  typ.Elem(),
	}
	_, exists := param.ResolveProvider(c)
	return exists
}

// Cleanup runs destructors in reverse order that was been created.
func (c *Container) Cleanup() {
	for i := len(c.cleanups) - 1; i >= 0; i-- {
		c.cleanups[i]()
	}
}

// struct that contains constructor with options.
type provideOptions struct {
	frame       stacktrace.Frame
	constructor Constructor
	options     []ProvideOption
}

// struct that contains invoke function with options.
type invokeOptions struct {
	frame   stacktrace.Frame
	fn      Invocation
	options []InvokeOption
}

// struct that container resolve target with options.
type resolveOptions struct {
	frame   stacktrace.Frame
	target  interface{}
	options []ResolveOption
}
