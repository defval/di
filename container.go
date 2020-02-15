package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/graph"
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
		graph:    graph.New(),
		cleanups: make([]func(), 0, 8),
	}
	// apply container options
	for _, opt := range options {
		opt.apply(c)
	}
	return c
}

// Container is a dependency injection container.
type Container struct {
	compiled bool             // compile state
	graph    *graph.Graph     // graph storage
	ctors    []provideOptions // initial provides
	invokes  []invokeOptions  // initial invocations
	resolves []resolveOptions // initial resolves
	cleanups []func()         // cleanup functions

}

// Provide adds constructor into container with parameters. It creates provider for constructor
// and place it into graph.
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
	// if provider already exists replace it to stub
	if c.graph.Exists(prov.ID()) {
		// duplicate types not allowed
		return fmt.Errorf("%s already exists in dependency graph", prov.ID())
	}
	// if prototype option provided wrap provider as singleton
	if !params.IsPrototype {
		prov = asSingleton(prov)
	}
	// add provider to graph
	c.graph.AddNode(providerNode{prov})
	// parse embed parameters
	for _, param := range prov.ParameterList() {
		if param.embed {
			c.graph.AddNode(providerNode{providerFromEmbedParameter(param)})
		}
	}
	iprovs := make([]*providerInterface, 0, len(params.Interfaces))
	// process interfaces
	for _, i := range params.Interfaces {
		// create interface provider
		iprov, err := newProviderInterface(prov, i)
		if err != nil {
			return err
		}
		if !c.graph.AddNode(providerNode{iprov}) {
			// if graph node already exists, resolve it as interface restricted
			// but in group may exists
			stub := newProviderStub(iprov.ID(), "have several implementations")
			c.graph.Replace(providerNode{stub})
		}
		iprovs = append(iprovs, iprov)
	}
	// process group for interfaces
	for _, iprov := range iprovs {
		group := newProviderGroup(iprov.ID())
		// if group node already exists use it
		if !c.graph.AddNode(providerNode{group}) {
			existing, err := c.graph.Node(group.ID())
			if err != nil {
				return err
			}
			group = existing.(providerNode).provider.(*providerGroup)
		}
		group.Add(iprov.provider.ID())
	}
	return nil
}

// Compile compiles the container. It iterates over all nodes
// in graph and register their parameters. Also container invoke functions provided
// by di.Invoke() container option.
func (c *Container) Compile(options ...CompileOption) error {
	// for _, opt := range options {
	// 	opt.apply(c)
	// }
	// process constructors
	for _, provide := range c.ctors {
		if err := c.Provide(provide.constructor, provide.options...); err != nil {
			return err
		}
	}
	// connect graph nodes, register provider parameters
	for _, node := range c.graph.Nodes() {
		provider := node.(providerNode)
		for _, param := range provider.ParameterList() {
			// node parameter provider
			pp, exists := param.ResolveProvider(c)
			if exists {
				if err := c.graph.AddEdge(pp.ID(), provider.ID(), 1); err != nil {
					return err
				}
				continue
			}
			if !exists && !param.optional {
				return fmt.Errorf("%s: dependency %s not exists in container", provider.ID(), param)
			}
		}
	}
	if err := graph.CheckCycles(c.graph); err != nil {
		return err
	}
	c.compiled = true
	// call initial invokes
	for _, fn := range c.invokes {
		if err := c.Invoke(fn.invocation, fn.options...); err != nil {
			return err
		}
	}
	// initial resolves
	for _, res := range c.resolves {
		if err := c.Resolve(res.target, res.options...); err != nil {
			return err
		}
	}
	return nil
}

// Resolve builds instance of target type and fills target pointer.
func (c *Container) Resolve(into interface{}, options ...ResolveOption) error {
	if !c.compiled {
		return fmt.Errorf("container not compiled")
	}
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
	if isEmbedParameter(typ) {
		return fmt.Errorf("resolve target must be a pointer, got di.Parameter")
	}
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
	if !c.compiled {
		return fmt.Errorf("container not compiled")
	}
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

// Exists checks that type exists in container, if not it cause error.
func (c *Container) Exists(target interface{}, options ...ResolveOption) bool {
	if !c.compiled {
		return false
	}
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
		name:  params.Name,
		typ:   typ.Elem(),
		embed: isEmbedParameter(typ.Elem()),
	}
	_, exists := param.ResolveProvider(c)
	return exists
}

// Cleanup runs destructors in order that was been created.
func (c *Container) Cleanup() {
	for i := len(c.cleanups) - 1; i >= 0; i-- {
		c.cleanups[i]()
	}
}

// struct that contains constructor with options.
type provideOptions struct {
	constructor Constructor
	options     []ProvideOption
}

// struct that contains invoke function with options.
type invokeOptions struct {
	invocation Invocation
	options    []InvokeOption
}

// struct that container resolve target with options.
type resolveOptions struct {
	target  interface{}
	options []ResolveOption
}
