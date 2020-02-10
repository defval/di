package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/graphkv"
	"github.com/goava/di/internal/reflection"
)

// New creates new container with provided options. Example usage:
//
// 	func NewHTTPServer(handler http.Handler) *http.Server {
//   	return &http.Server{}
// 	}
//
// 	func NewHTTPServeMux() *http.ServeMux {
//   	return http.ServeMux{}
// 	}
//
// 	container := di.New(
//   	di.Provide(NewHTTPServer),
//   	di.Provide(NewHTTPServeMux, di.As()),
// 	)
func New(options ...Option) *Container {
	c := &Container{
		compiled: false,
		graph:    graphkv.New(),
		cleanups: make([]func(), 0, 8),
	}
	for _, opt := range options {
		opt.apply(c)
	}
	return c
}

// Container is a dependency injection container.
type Container struct {
	compiled    bool
	graph       *graphkv.Graph
	cleanups    []func()
	provides    []constructorOptions
	invocations []invocationOptions
}

// Provide adds constructor into container with parameters.
func (c *Container) Provide(constructor interface{}, options ...ProvideOption) {
	params := ProvideParams{}
	for _, opt := range options {
		opt.apply(&params)
	}
	provider := internalProvider(newProviderConstructor(params.Name, constructor))
	key := provider.Key()
	if c.graph.Exists(key) {
		panicf("The `%s` type already exists in container", provider.Key())
	}
	if !params.IsPrototype {
		provider = asSingleton(provider)
	}
	// add provider to graph
	c.graph.Add(key, provider)
	// parse embed parameters
	for _, param := range provider.ParameterList() {
		if param.embed {
			embed := newProviderEmbed(param)
			c.graph.Add(embed.Key(), embed)
		}
	}
	// provide parameter bag
	if len(params.Parameters) != 0 {
		parameterBugProvider := createParameterBugProvider(provider.Key(), params.Parameters)
		c.graph.Add(parameterBugProvider.Key(), parameterBugProvider)
	}
	// process interfaces
	for _, iface := range params.Interfaces {
		c.processProviderInterface(provider, iface)
	}
}

// Compile compiles the container. It iterates over all nodes
// in graph and register their parameters.
func (c *Container) Compile() {
	for _, provide := range c.provides {
		// todo: add error
		c.Provide(provide.constructor, provide.options...)
	}
	graphProvider := func() *Graph { return &Graph{graph: c.graph.DOTGraph()} }
	c.Provide(graphProvider)
	for _, node := range c.graph.Nodes() {
		c.registerProviderParameters(node.Value.(internalProvider))
	}
	if err := c.graph.CheckCycles(); err != nil {
		panic(err.Error())
	}
	// call invocations
	for _, fn := range c.invocations {
		if err := c.Invoke(fn.invocation, fn.options...); err != nil {
			panic(err.Error()) // todo: remove panic
		}
	}
	c.compiled = true
}

// Resolve builds instance of target type and fills target pointer.
func (c *Container) Resolve(into interface{}, options ...ExtractOption) error {
	if !c.compiled {
		return fmt.Errorf("container not compiled")
	}
	if into == nil {
		return fmt.Errorf("resolve target must be a pointer, got `nil`")
	}
	if !reflection.IsPtr(into) {
		return fmt.Errorf("resolve target must be a pointer, got `%s`", reflect.TypeOf(into))
	}
	params := ExtractParams{}
	for _, opt := range options {
		opt.apply(&params)
	}
	typ := reflect.TypeOf(into)
	param := parameter{
		name:  params.Name,
		res:   typ.Elem(),
		embed: isEmbedParameter(typ),
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
func (c *Container) Invoke(fn interface{}, options ...InvokeOption) error {
	if !c.compiled {
		return fmt.Errorf("container not compiled")
	}
	params := InvokeParams{}
	for _, opt := range options {
		opt.apply(&params)
	}
	invoker, err := newInvoker(fn)
	if err != nil {
		return err
	}
	return invoker.Invoke(c)
}

// Exists
func (c *Container) Exists(target interface{}, options ...ExtractOption) bool {
	if !c.compiled {
		return false
	}
	if target == nil {
		return false
	}
	params := ExtractParams{}
	for _, opt := range options {
		opt.apply(&params)
	}
	typ := reflect.TypeOf(target)
	param := parameter{
		name:  params.Name,
		res:   typ.Elem(),
		embed: isEmbedParameter(typ),
	}
	_, exists := param.ResolveProvider(c)
	return exists
}

// Cleanup runs destructors in order that was been created.
func (c *Container) Cleanup() {
	for _, cleanup := range c.cleanups {
		cleanup()
	}
}

// processProviderInterface represents instances as interfaces and groups.
func (c *Container) processProviderInterface(provider internalProvider, as interface{}) {
	// create interface from provider
	iface := newProviderInterface(provider, as)
	key := iface.Key()
	if c.graph.Exists(key) {
		stub := newProviderStub(key, "have several implementations")
		c.graph.Replace(key, stub)
	} else {
		// add interface node
		c.graph.Add(key, iface)
	}
	// create group
	group := newProviderGroup(key)
	groupKey := group.Key()
	// check exists
	if c.graph.Exists(groupKey) {
		// if exists use existing group
		node := c.graph.Get(groupKey)
		group = node.Value.(*providerGroup)
	} else {
		// else add new group to graph
		c.graph.Add(groupKey, group)
	}
	// add provider reference into group
	providerKey := provider.Key()
	group.Add(providerKey)
}

// registerProviderParameters registers provider parameters in a dependency graph.
func (c *Container) registerProviderParameters(p internalProvider) {
	for _, param := range p.ParameterList() {
		provider, exists := param.ResolveProvider(c)
		if exists {
			c.graph.Edge(provider.Key(), p.Key())
			continue
		}
		if !exists && !param.optional {
			panicf("%s: dependency %s not exists in container", p.Key(), param)
		}
	}
}

// struct that contains constructor with options.
type constructorOptions struct {
	constructor Constructor
	options     []ProvideOption
}

// struct that contains invocation with options.
type invocationOptions struct {
	invocation Invocation
	options    []InvokeOption
}
