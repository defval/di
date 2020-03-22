package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/graph"
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
		compiled: false,
		graph:    graph.New(),
		cleanups: make([]func(), 0, 8),
		logger:   nopLogger{},
	}
	// apply container options
	for _, opt := range options {
		opt.apply(c)
	}
	// process constructors
	provideErr := errProvideFailed{}
	for _, provide := range c.provides {
		if err := c.Provide(provide.constructor, provide.options...); err != nil {
			provideErr = provideErr.Append(provide.frame, err)
		}
	}
	if len(provideErr) > 0 {
		return nil, provideErr
	}
	if !c.mcf {
		if err := c.Compile(); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// Container is a dependency injection container.
type Container struct {
	mcf      bool             // manual compile flow - if true compile will not be called on New()
	compiled bool             // compile state
	graph    *graph.Graph     // graph storage
	provides []provideOptions // initial provides
	invokes  []invokeOptions  // initial invokes
	resolves []resolveOptions // initial resolves
	cleanups []func()         // cleanup functions
	logger   Logger           // internal logger
}

// Provide provides to container reliable way to build type. The constructor will be invoked lazily on-demand.
// For more information about constructors see Constructor interface. ProvideOption can add additional behavior to
// the process of type resolving.
func (c *Container) Provide(constructor Constructor, options ...ProvideOption) (err error) {
	if c.compiled {
		return fmt.Errorf("dependency providing restricted after container compile")
	}
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
			// if node add returns false graph node already exists
			// error can be omitted
			existing, _ := c.graph.Node(group.ID())
			group = existing.(providerNode).provider.(*providerGroup)
		}
		group.Add(iprov.provider.ID())
	}
	return nil
}

// Compile compiles the container. First, it iterates over all definitions and register their
// parameters. Container links definitions with each other and checks that result dependency graph is not cyclic.
// In final, the container invoke functions provided by di.Invoke() container option and resolves types
// provided by di.Resolve() container option. Between invokes and resolves, the container tries to find di.Logger
// interface, and if it is found sets it as an internal logger.
func (c *Container) Compile(_ ...CompileOption) error {
	if c.compiled {
		return fmt.Errorf("container already compiled, recompilation restricted")
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
	// error omitted because if logger could not be resolve it will be default
	_ = c.Resolve(&c.logger)
	// call initial invokes
	for _, invoke := range c.invokes {
		if err := c.Invoke(invoke.fn, invoke.options...); err != nil {
			return errInvokeFailed{invoke.frame, err}
		}
	}
	// initial resolves
	for _, resolve := range c.resolves {
		if err := c.Resolve(resolve.target, resolve.options...); err != nil {
			return errResolveFailed{resolve.frame, err}
		}
	}
	return nil
}

// Resolve builds instance of target type and fills target pointer.
func (c *Container) Resolve(into interface{}, options ...ResolveOption) error {
	if !c.compiled {
		return fmt.Errorf("container not compiled, resolve dependencies possible only after compilation")
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
		return fmt.Errorf("container not compiled, function invokes possible only after compilation")
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
		name:  params.Name,
		typ:   typ.Elem(),
		embed: isEmbedParameter(typ.Elem()),
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

// Parameter is embed helper that indicates that type is a constructor embed parameter.
type Parameter struct {
	internalParameter
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
