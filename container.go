package di

import (
	"fmt"
	"reflect"

	"github.com/goava/di/internal/reflection"
)

// Container is a dependency injection container.
type Container struct {
	// Logger that logs internal actions.
	logger Logger
	// Initial options will be processed on di.New().
	initial struct {
		// Array of di.Provide() options.
		provides []provideOptions
		// Array of di.Invoke() options.
		invokes []invokeOptions
		// Array of di.Resolve() options.
		resolves []resolveOptions
	}
	// Mapping from key to provider that can provide value for that key.
	providers map[reflect.Type]*providerList
	// Array of provider cleanups.
	cleanups []func()
	// Flag indicates acyclic verification state
	verified map[key]bool
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
		logger:    nopLogger{},
		providers: map[reflect.Type]*providerList{},
		cleanups:  []func(){},
		verified:  map[key]bool{},
	}
	// apply container options
	for _, opt := range options {
		opt.apply(c)
	}
	// process di.Provide() options
	for _, provide := range c.initial.provides {
		if err := c.provide(provide.constructor, provide.options...); err != nil {
			return nil, ErrProvideFailed{
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
		if err != nil && isUsageError(err) {
			return nil, ErrInvokeFailed{invoke.frame, err}
		}
		if err != nil {
			return nil, err
		}
	}
	// process di.Resolve() options
	for _, resolve := range c.initial.resolves {
		if err := c.resolve(resolve.target, resolve.options...); err != nil {
			return nil, ErrResolveFailed{resolve.frame, err}
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
	fn, valid := reflection.InspectFunc(constructor)
	if !valid {
		return fmt.Errorf("invalid constructor signature, got %s", reflect.TypeOf(constructor))
	}
	params := ProvideParams{}
	// apply provide options
	for _, opt := range options {
		opt.apply(&params)
	}
	// create constructor provider
	p, err := newProviderConstructor(params.Name, fn)
	if err != nil {
		return err
	}
	cleanup := p.ctorType == ctorCleanup || p.ctorType == ctorCleanupError
	if cleanup && params.IsPrototype {
		return fmt.Errorf("cleanup not supported with prototype providers")
	}
	// provider list
	plist, ok := c.providers[p.Type()]
	if !ok {
		// create list of providers
		plist = createProviderList()
		c.providers[p.Type()] = plist
	}
	fp := provider(p)
	if !params.IsPrototype {
		fp = asSingleton(p)
	}
	uniq, err := plist.Add(fp)
	if err != nil {
		return err
	}
	link := keyUniq{key{fp.Type(), fp.Name()}, uniq}
	if err := c.processInterfaces(link, params.Interfaces, params.IsPrototype); err != nil {
		return err
	}
	return nil
}

func (c *Container) processInterfaces(key keyUniq, interfaces []Interface, isPrototype bool) error {
	// interface raw
	for _, iraw := range interfaces {
		// provider interface
		piface, err := newProviderInterface(key.uniq, key.key, iraw)
		if err != nil {
			return err
		}
		// interface list
		ilist, ok := c.providers[piface.Type()]
		if !ok {
			ilist = createProviderList()
			c.providers[piface.Type()] = ilist

		}
		fpiface := provider(piface)
		if !isPrototype {
			fpiface = asSingleton(piface)
		}
		_, err = ilist.Add(fpiface)
		if err != nil {
			return err
		}
	}
	return nil
}

// Resolve resolves type and fills target pointer.
//
//	var server *http.Server
//	if err := container.Resolve(&server); err != nil {
//		// handle error
//	}
func (c *Container) Resolve(into interface{}, options ...ResolveOption) error {
	if err := c.resolve(into, options...); err != nil {
		return resolveErrWithStack(err)
	}
	return nil
}

func (c *Container) resolve(into interface{}, options ...ResolveOption) error {
	if into == nil {
		return fmt.Errorf("resolve target must be a pointer, got nil")
	}
	if reflect.ValueOf(into).Kind() != reflect.Ptr {
		return fmt.Errorf("resolve target must be a pointer, got %s", reflect.TypeOf(into))
	}
	params := ResolveParams{}
	// apply extract options
	for _, opt := range options {
		opt.apply(&params)
	}
	param := parameter{
		name: params.Name,
		typ:  reflect.TypeOf(into).Elem(),
	}
	// check cycle verified
	if !c.verified[param.Key()] {
		err := checkCycles(c, param)
		if err != nil {
			return err
		}
		c.verified[param.Key()] = true
	}
	value, err := param.ResolveValue(c)
	if err != nil {
		return err
	}
	targetValue := reflect.ValueOf(into).Elem()
	targetValue.Set(value)
	return nil
}

// Invoke calls the function fn. It parses function parameters. Looks for it in a container.
// And invokes function with them. See Invocation for details.
func (c *Container) Invoke(invocation Invocation, options ...InvokeOption) error {
	err := c.invoke(invocation, options...)
	if err != nil && isUsageError(err) {
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
	fn, isFn := reflection.InspectFunc(invocation)
	if !isFn {
		return errInvalidInvocation{fmt.Errorf("invalid invocation signature, got %s", fn.Type)}
	}
	if !validateInvocation(fn) {
		return errInvalidInvocation{fmt.Errorf("invalid invocation signature, got %s", fn.Type)}
	}
	var plist parameterList
	for j := 0; j < fn.NumIn(); j++ {
		param := parameter{
			typ:      fn.In(j),
			optional: false,
		}
		// check cycle verified
		if !c.verified[param.Key()] {
			err := checkCycles(c, param)
			if err != nil {
				return errInvalidInvocation{err}
			}
			c.verified[param.Key()] = true
		}
		plist = append(plist, param)
	}
	values, err := plist.Resolve(c)
	if err != nil {
		return err
	}
	results := reflection.CallResult(fn.Call(values))
	if len(results) == 0 {
		return nil
	}
	return results.Error(0)
}

// Has checks that type exists in container, if not it return false.
//
// 	var server *http.Server
//	if container.Has(&server) {
//		// handle server existence
//	}
//
// It like Resolve() but doesn't instantiate a type.
func (c *Container) Has(into interface{}, options ...ResolveOption) bool {
	if into == nil {
		return false
	}
	params := ResolveParams{}
	// apply options
	for _, opt := range options {
		opt.apply(&params)
	}
	typ := reflect.TypeOf(into)
	param := parameter{
		name: params.Name,
		typ:  typ.Elem(),
	}
	_, err := param.ResolveProvider(c)
	if err == nil {
		return true
	}
	if _, ok := err.(errParameterProviderNotFound); ok {
		return false
	}
	bug()
	return false
}

// Cleanup runs destructors in reverse order that was been created.
func (c *Container) Cleanup() {
	for i := len(c.cleanups) - 1; i >= 0; i-- {
		c.cleanups[i]()
	}
}

// Compile compiles the container.
// Deprecated: Compile deprecated: https://github.com/goava/di/pull/9
func (c *Container) Compile(_ ...CompileOption) error { return nil }
