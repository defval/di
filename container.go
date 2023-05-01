package di

import (
	"errors"
	"fmt"
	"reflect"
)

// Container is a dependency injection container.
type Container struct {
	// Dependency injection schema.
	schema *defaultSchema
	// Array of provider cleanups.
	cleanups []func()
}

// New constructs container with provided options. Example usage (simplified):
//
// Define constructors and invocations:
//
//	func NewHTTPServer(mux *http.ServeMux) *http.Server {
//		return &http.Server{
//			Handler: mux,
//		}
//	}
//
//	func NewHTTPServeMux() *http.ServeMux {
//		return http.ServeMux{}
//	}
//
//	func StartServer(server *http.Server) error {
//		return server.ListenAndServe()
//	}
//
// Use it with container:
//
//	container, err := di.New(
//		di.Provide(NewHTTPServer),
//		di.Provide(NewHTTPServeMux),
//		di.Invoke(StartServer),
//	)
//	if err != nil {
//		// handle error
//	}
func New(options ...Option) (_ *Container, err error) {
	c := &Container{
		schema:   newDefaultSchema(),
		cleanups: []func(){},
	}
	var di diopts
	// apply container diopts
	for _, opt := range options {
		opt.apply(&di)
	}
	// provide container to advanced usage e.g. condition providing
	_ = c.provide(func() *Container { return c })
	if err := c.apply(di); err != nil {
		return nil, err
	}
	return c, nil
}

// Apply applies options to container.
//
//	err := container.Apply(
//		di.Provide(NewHTTPServer),
//	)
//	if err != nil {
//		// handle error
//	}
func (c *Container) Apply(options ...Option) error {
	var di diopts
	for _, opt := range options {
		opt.apply(&di)
	}
	return c.apply(di)
}

// Provide provides to container reliable way to build type. The constructor will be invoked lazily on-demand.
// For more information about constructors see Constructor interface. ProvideOption can add additional behavior to
// the process of type resolving.
func (c *Container) Provide(constructor Constructor, options ...ProvideOption) error {
	if err := c.provide(constructor, options...); err != nil {
		return errWithStack(err)
	}
	return nil
}

// ProvideValue provides value as is.
func (c *Container) ProvideValue(value Value, options ...ProvideOption) error {
	if err := c.provideValue(value, options...); err != nil {
		return errWithStack(err)
	}
	return nil
}

// Invocation is a function whose signature looks like:
//
//	func StartServer(server *http.Server) error {
//		return server.ListenAndServe()
//	}
//
// Like a constructor invocation may have unlimited count of arguments and
// they will be resolved automatically. The invocation can return an optional error.
// Error will be returned as is.
type Invocation interface{}

// Invoke calls the function fn. It parses function parameters. Looks for it in a container.
// And invokes function with them. See Invocation for details.
func (c *Container) Invoke(invocation Invocation, options ...InvokeOption) error {
	err := c.invoke(invocation, options...)
	if err != nil && knownError(err) {
		return errWithStack(err)
	}
	if err != nil {
		return err
	}
	return nil
}

type Pointer interface{}

// Has checks that type exists in container, if not it return false.
//
//	var server *http.Server
//	if container.Has(&server) {
//		// handle server existence
//	}
//
// It like Resolve() but doesn't instantiate a type.
func (c *Container) Has(target Pointer, options ...ResolveOption) (bool, error) {
	if _, err := c.find(target, options...); errors.Is(err, ErrTypeNotExists) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Resolve resolves type and fills target pointer.
//
//	var server *http.Server
//	if err := container.Resolve(&server); err != nil {
//		// handle error
//	}
func (c *Container) Resolve(ptr Pointer, options ...ResolveOption) error {
	if err := c.resolve(ptr, options...); err != nil {
		return errWithStack(err)
	}
	return nil
}

// ValueFunc is a lazy-loading wrapper for iteration.
type ValueFunc func() (interface{}, error)

// IterateFunc function that will be called on each instance in iterate selection.
type IterateFunc func(tags Tags, value ValueFunc) error

// Iterate iterates over group of Pointer type with IterateFunc.
//
//	 var servers []*http.Server
//	 iterFn := func(tags di.Tags, loader ValueFunc) error {
//			i, err := loader()
//			if err != nil {
//				return err
//			}
//			// do stuff with result: i.(*http.Server)
//			return nil
//	 }
//	 container.Iterate(&servers, iterFn)
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

// Cleanup runs destructors in reverse order that was been created.
func (c *Container) Cleanup() {
	for i := len(c.schema.cleanups) - 1; i >= 0; i-- {
		c.schema.cleanups[i]()
	}
}

// AddParent adds a parent container. Types are resolved from the container,
// it's parents, and ancestors. An error is a cycle is detected in ancestry tree.
func (c *Container) AddParent(parent *Container) error {
	return c.schema.addParent(parent.schema)
}

func (c *Container) apply(di diopts) error {
	for _, provide := range di.values {
		if err := c.provideValue(provide.value, provide.options...); err != nil {
			return fmt.Errorf("%s: %w", provide.frame, err)
		}
	}
	// process di.Resolve() diopts
	for _, provide := range di.provides {
		if err := c.provide(provide.constructor, provide.options...); err != nil {
			return fmt.Errorf("%s: %w", provide.frame, err)
		}
	}
	// error omitted because if logger could not be resolved it will be default
	// process di.Invoke() diopts
	for _, invoke := range di.invokes {
		err := c.invoke(invoke.fn, invoke.options...)
		if err != nil && knownError(err) {
			return fmt.Errorf("%s: %w", invoke.frame, err)
		}
		if err != nil {
			return err
		}
	}
	// process di.Resolve() diopts
	for _, resolve := range di.resolves {
		if err := c.resolve(resolve.target, resolve.options...); err != nil {
			return fmt.Errorf("%s: %w", resolve.frame, err)
		}
	}
	return nil
}

func (c *Container) provide(constructor Constructor, options ...ProvideOption) error {
	if constructor == nil {
		return fmt.Errorf("invalid constructor signature, got nil")
	}
	params := ProvideParams{}
	// apply provide options
	for _, opt := range options {
		opt.applyProvide(&params)
	}
	n, err := newConstructorNode(constructor)
	if err != nil {
		return err
	}
	n.decorators = params.Decorators
	for k, v := range params.Tags {
		n.tags[k] = v
	}
	return c.provideNode(n, params)
}

func (c *Container) provideValue(value Value, options ...ProvideOption) error {
	if value == nil {
		return fmt.Errorf("invalid value, got nil")
	}
	params := ProvideParams{}
	// apply provide diopts
	for _, opt := range options {
		opt.applyProvide(&params)
	}
	v := reflect.ValueOf(value)
	n := &node{
		compiler: valueCompiler{
			rv: v,
		},
		rv:         new(reflect.Value),
		rt:         v.Type(),
		tags:       params.Tags,
		decorators: params.Decorators,
	}
	return c.provideNode(n, params)
}

func (c *Container) provideNode(n *node, params ProvideParams) error {
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
			rv:         n.rv,
			rt:         i.Type,
			tags:       n.tags,
			compiler:   n.compiler,
			decorators: n.decorators,
		})
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
		return fmt.Errorf("%s: %w", node, err)
	}
	rv := reflect.ValueOf(ptr)
	target := rv.Elem()
	if canInject(rv.Type()) {
		for index := range parsePopulateFields(target.Type()) {
			target.Field(index).Set(value.Field(index))
		}
	} else {
		target.Set(value)
	}
	return nil
}

func (c *Container) invoke(invocation Invocation, _ ...InvokeOption) error {
	// params := InvokeParams{}
	// for _, opt := range diopts {
	// 	opt.apply(&params)
	// }
	if invocation == nil {
		return fmt.Errorf("%w, got %s", errInvalidInvocationSignature, "nil")
	}
	fn, valid := inspectFunction(invocation)
	if !valid {
		return fmt.Errorf("%w, got %s", errInvalidInvocationSignature, reflect.TypeOf(invocation))
	}
	if !validateInvocation(fn) {
		return fmt.Errorf("%w, got %s", errInvalidInvocationSignature, reflect.TypeOf(invocation))
	}
	nodes, err := parseInvocationParameters(fn, c.schema)
	if err != nil {
		return err
	}
	var args []reflect.Value
	for _, node := range nodes {
		if err := c.schema.prepare(node); err != nil {
			return err
		}
		v, err := node.Value(c.schema)
		if err != nil {
			return fmt.Errorf("%s: %s", node, err)
		}
		args = append(args, v)
	}
	res := funcResult(fn.Call(args))
	if len(res) == 0 {
		return nil
	}
	return res.error(0)
}

func (c *Container) find(ptr Pointer, options ...ResolveOption) (*node, error) {
	if ptr == nil {
		return nil, fmt.Errorf("target must be a pointer, got nil")
	}
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("target must be a pointer, got %s", reflect.TypeOf(ptr))
	}
	params := ResolveParams{}
	// apply extract diopts
	for _, opt := range options {
		opt.applyResolve(&params)
	}
	node, err := c.schema.find(reflect.TypeOf(ptr).Elem(), params.Tags)
	if err != nil {
		return nil, err
	}
	if err := c.schema.prepare(node); err != nil {
		return nil, err
	}
	return node, nil
}

type diopts struct {
	// Array of di.Provide() options.
	provides []provideOptions
	// Array of di.ProvideValue() options.
	values []provideValueOptions
	// Array of di.Invoke() options.
	invokes []invokeOptions
	// Array of di.Resolve() options.
	resolves []resolveOptions
}
