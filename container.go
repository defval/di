package di

import (
	"fmt"
	"reflect"
)

// Container is a dependency injection container.
type Container struct {
	// Logger that logs internal actions.
	logger Logger
	// Dependency injection schema.
	schema *defaultSchema
	// Initial options will be processed on di.New().
	initial struct {
		// Array of di.Provide() options.
		provides []provideOptions
		// Array of di.Invoke() options.
		invokes []invokeOptions
		// Array of di.Resolve() options.
		resolves []resolveOptions
	}
	// Array of provider cleanups.
	cleanups []func()
	// Flag indicates acyclic verification state
	//verified map[key]bool
}

// New constructs container with provided options. Example usage (simplified):
//
// Define constructors and invocations:
//
// 	func NewHTTPServer(mux *http.ServeMux) *http.Server {
// 		return &http.Server{
// 			Handler: mux,
// 		}
// 	}
//
// 	func NewHTTPServeMux() *http.ServeMux {
// 		return http.ServeMux{}
// 	}
//
// 	func StartServer(server *http.Server) error {
//		return server.ListenAndServe()
//	}
//
// Use it with container:
//
// 	container, err := di.New(
// 		di.Provide(NewHTTPServer),
// 		di.Provide(NewHTTPServeMux),
//		di.Invoke(StartServer),
// 	)
// 	if err != nil {
//		// handle error
//	}
func New(options ...Option) (_ *Container, err error) {
	c := &Container{
		logger:   nopLogger{},
		schema:   newDefaultSchema(),
		cleanups: []func(){},
	}
	// apply container options
	for _, opt := range options {
		opt.apply(c)
	}
	// process di.Resolve() options
	for _, provide := range c.initial.provides {
		if err := c.provide(provide.constructor, provide.options...); err != nil {
			return nil, errProvideFailed{
				provide.frame,
				err,
			}
		}
	}
	// provide container to advanced usage e.g. condition providing
	_ = c.provide(func() *Container { return c })
	// error omitted because if logger could not be resolved it will be default
	_ = c.resolve(&c.logger)
	// process di.Invoke() options
	for _, invoke := range c.initial.invokes {
		err := c.invoke(invoke.fn, invoke.options...)
		if err != nil && knownError(err) {
			return nil, errInvokeFailed{invoke.frame, err}
		}
		if err != nil {
			return nil, err
		}
	}
	// process di.Resolve() options
	for _, resolve := range c.initial.resolves {
		if err := c.resolve(resolve.target, resolve.options...); err != nil {
			return nil, errResolveFailed{resolve.frame, err}
		}
	}
	return c, nil
}

// Provide provides to container reliable way to build type. The constructor will be invoked lazily on-demand.
// For more information about constructors see Constructor interface. ProvideOption can add additional behavior to
// the process of type resolving.
func (c *Container) Provide(constructor Constructor, options ...ProvideOption) error {
	if err := c.provide(constructor, options...); err != nil {
		return provideErrWithStack(err)
	}
	return nil
}

func (c *Container) provide(constructor Constructor, options ...ProvideOption) error {
	if constructor == nil {
		return fmt.Errorf("invalid constructor signature, got nil")
	}
	fn, valid := inspectFunction(constructor)
	if !valid {
		return fmt.Errorf("invalid constructor signature, got %s", reflect.TypeOf(constructor))
	}
	params := ProvideParams{}
	// apply provide options
	for _, opt := range options {
		opt.applyProvide(&params)
	}
	n, err := nodeFromFunction(fn)
	if err != nil {
		return err
	}
	for k, v := range params.Tags {
		n.tags[k] = v
	}
	c.schema.register(n)
	// register interfaces
	for _, cur := range params.Interfaces {
		i, err := inspectInterfacePointer(cur)
		if err != nil {
			return err
		}
		if !n.rt.Implements(i.Type) {
			return fmt.Errorf("%s not implement %s", n, i.Type)
		}
		c.schema.register(&node{
			rv:       n.rv,
			rt:       i.Type,
			tags:     n.tags,
			compiler: n.compiler,
		})
		return nil
	}
	return nil
}

type Pointer interface{}

// Resolve resolves type and fills target pointer.
//
//	var server *http.Server
//	if err := container.Resolve(&server); err != nil {
//		// handle error
//	}
func (c *Container) Resolve(ptr Pointer, options ...ResolveOption) error {
	if err := c.resolve(ptr, options...); err != nil {
		return resolveErrWithStack(err)
	}
	return nil
}

func (c *Container) resolve(ptr Pointer, options ...ResolveOption) error {
	node, err := c.find(ptr, options...)
	if err != nil {
		return err
	}
	value, err := node.Value(c.schema)
	if err != nil {
		return fmt.Errorf("%s: %s", node, err)
	}
	target := reflect.ValueOf(ptr).Elem()
	target.Set(value)
	return nil
}

// Invoke calls the function fn. It parses function parameters. Looks for it in a container.
// And invokes function with them. See Invocation for details.
func (c *Container) Invoke(invocation Invocation, options ...InvokeOption) error {
	err := c.invoke(invocation, options...)
	if err != nil && knownError(err) {
		return invokeErrWithStack(err)
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *Container) invoke(invocation Invocation, _ ...InvokeOption) error {
	// params := InvokeParams{}
	// for _, opt := range options {
	// 	opt.apply(&params)
	// }
	if invocation == nil {
		return errInvalidInvocation{fmt.Errorf("invalid invocation signature, got %s", "nil")}
	}
	fn, valid := inspectFunction(invocation)
	if !valid {
		return errInvalidInvocation{fmt.Errorf("invalid invocation signature, got %s", fn.Type)}
	}
	if !validateInvocation(fn) {
		return errInvalidInvocation{fmt.Errorf("invalid invocation signature, got %s", fn.Type)}
	}
	nodes, err := getFunctionNodes(fn, c.schema)
	if err != nil {
		return errInvalidInvocation{err}
	}
	var args []reflect.Value
	for _, node := range nodes {
		if err := prepare(c.schema, node); err != nil {
			return errInvalidInvocation{err}
		}
		v, err := node.Value(c.schema)
		if err != nil {
			return errInvalidInvocation{err}
		}
		args = append(args, v)
	}
	res := funcOut(fn.Call(args))
	if len(res) == 0 {
		return nil
	}
	return res.error(0)
}

// Has checks that type exists in container, if not it return false.
//
// 	var server *http.Server
//	if container.Has(&server) {
//		// handle server existence
//	}
//
// It like Resolve() but doesn't instantiate a type.
func (c *Container) Has(target Pointer, options ...ResolveOption) bool {
	if _, err := c.find(target, options...); err != nil {
		return false
	}
	return true
}

// ValueFunc is a lazy-loading wrapper for iteration.
type ValueFunc func() (interface{}, error)

// IterateFunc function that will be called on each instance in iterate selection.
type IterateFunc func(tags Tags, value ValueFunc) error

// Iterate iterates over group of Pointer type with IterateFunc.
//
//  var servers []*http.Server
//  iterFn := func(tags di.Tags, loader ValueFunc) error {
//		i, err := loader()
//		if err != nil {
//			return err
//		}
//		// do stuff with result: i.(*http.Server)
//		return nil
//  }
//  container.Iterate(&servers, iterFn)
func (c *Container) Iterate(target Pointer, fn IterateFunc, options ...ResolveOption) error {
	node, err := c.find(target, options...)
	if err != nil {
		return err
	}
	group, ok := node.compiler.(*groupCompiler)
	if ok {
		for i, n := range group.matched {
			err = fn(n.tags, func() (interface{}, error) {
				v, err := n.Value(c.schema)
				if err != nil {
					return nil, err
				}
				return v.Interface(), nil
			})
			if err != nil {
				return fmt.Errorf("%s with index %d failed: %s", node, i, err)
			}
		}
		return nil
	}
	return fmt.Errorf("iteration can be used with groups only")
}

func (c *Container) find(ptr Pointer, options ...ResolveOption) (*node, error) {
	if ptr == nil {
		return nil, fmt.Errorf("target must be a pointer, got nil")
	}
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("target must be a pointer, got %s", reflect.TypeOf(ptr))
	}
	params := ResolveParams{}
	// apply extract options
	for _, opt := range options {
		opt.applyResolve(&params)
	}
	node, err := c.schema.find(reflect.TypeOf(ptr).Elem(), params.Tags)
	if err != nil {
		return nil, err
	}
	if err := prepare(c.schema, node); err != nil {
		return nil, err
	}
	return node, nil
}

// Cleanup runs destructors in reverse order that was been created.
func (c *Container) Cleanup() {
	for i := len(c.schema.cleanups) - 1; i >= 0; i-- {
		c.schema.cleanups[i]()
	}
}
