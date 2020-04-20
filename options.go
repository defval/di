package di

import (
	"github.com/goava/di/internal/stacktrace"
)

// Option is a functional option that configures container. If you don't know about functional
// options, see https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis.
// Below presented all possible options with their description:
//
// 	- di.Provide - provide constructors
//	- di.Invoke - add invocations
//	- di.Resolve - resolves type
type Option interface {
	apply(c *Container)
}

// Provide returns container option that provides to container reliable way to build type. The constructor will
// be invoked lazily on-demand. For more information about constructors see Constructor interface. ProvideOption can
// add additional behavior to the process of type resolving.
func Provide(constructor Constructor, options ...ProvideOption) Option {
	frame := stacktrace.CallerFrame(0)
	return option(func(c *Container) {
		c.initial.provides = append(c.initial.provides, provideOptions{
			frame,
			constructor,
			options,
		})
	})
}

// Constructor is a function with follow signature:
//
// 	func NewHTTPServer(addr string, handler http.Handler) (server *http.Server, cleanup func(), err error) {
// 		server := &http.Server{
// 			Addr: addr,
// 		}
// 		cleanup = func() {
// 			server.Close()
// 		}
// 		return server, cleanup, nil
// 	}
//
// This constructor function teaches container how to build server. Arguments (addr and handler) in this function
// is a dependencies. They will be resolved automatically when someone needs a server. Constructor may have unlimited
// count of dependencies, but note that container should know how build each of them.
// Second result of this function is a optional cleanup callback. It describes that container will do on shutdown.
// Third result is a optional error. Sometimes our types cannot be constructed.
type Constructor interface{}

// ProvideOption is a functional option interface that modify provide behaviour. See di.As(), di.WithName().
type ProvideOption interface {
	apply(params *ProvideParams)
}

// As returns provide option that specifies interfaces for constructor resultant type.
//
// INTERFACE USAGE:
//
// You can provide type as interface and resolve it later without using of direct implementation.
// This creates less cohesion of code and promotes be more testable.
//
// Create type constructors:
//
// 		func NewServeMux() *http.ServeMux {
// 			return &http.ServeMux{}
// 		}
//
//		func NewServer(handler *http.Handler) *http.Server {
//			return &http.Server{
//				Handler: handler,
//			}
//		}
//
// Build container with di.As provide option:
//
//		container, err := di.New(
//			di.Provide(NewServer),
//			di.Provide(NewServeMux, di.As(new(http.Handler)),
//		)
//		if err != nil {
//			// handle error
//		}
//		var server *http.Server
//		if err := container.Resolve(&http.Server); err != nil {
//			// handle error
//		}
//
// In this example you can see how container inject type *http.ServeMux as http.Handler
// interface into the server constructor.
//
// GROUP USAGE:
//
// Container automatically creates group for interfaces. For example, you can use type []http.Handler in
// previous example.
//
//		var handlers []http.Handler
//		if err := container.Resolve(&handlers); err != nil {
//			// handle error
//		}
//
// Container checks that provided type implements interface if not cause compile error.
func As(interfaces ...Interface) ProvideOption {
	return provideOption(func(params *ProvideParams) {
		params.Interfaces = append(params.Interfaces, interfaces...)
	})
}

// Interface is a pointer to interface, like new(http.Handler). Tell container that provided
// type may be used as interface.
type Interface interface{}

// WithName modifies Provide() behavior. It adds name identity for provided type.
func WithName(name string) ProvideOption {
	return provideOption(func(params *ProvideParams) {
		params.Name = name
	})
}

// Prototype modifies Provide() behavior. By default, each type resolves as a singleton. This option sets that
// each type resolving creates a new instance of the type.
//
//		container, err := di.New(
// 			Provide(NewHTTPServer, inject.Prototype())
//		)
//		if err != nil {
//			// handle error
//		}
// 		var server1, server2 *http.Server
// 		if err := container.Resolve(&server1); err != nil {
//			// handle error
//		}
//		if err := container.Resolve(&server2); err != nil {
//			// handle error
//		}
//
func Prototype() ProvideOption {
	return provideOption(func(params *ProvideParams) {
		params.IsPrototype = true
	})
}

// Resolve returns container options that resolves type into target. All resolves will be done on compile stage
// after call invokes.
func Resolve(target interface{}, options ...ResolveOption) Option {
	frame := stacktrace.CallerFrame(0)
	return option(func(c *Container) {
		c.initial.resolves = append(c.initial.resolves, resolveOptions{
			frame,
			target,
			options,
		})
	})
}

// Invoke returns container option that registers container invocation. All invocations
// will be called on di.New() after processing di.Provide() options.
// See Container.Invoke() for details.
func Invoke(fn Invocation, options ...InvokeOption) Option {
	frame := stacktrace.CallerFrame(0)
	return option(func(c *Container) {
		c.initial.invokes = append(c.initial.invokes, invokeOptions{
			frame,
			fn,
			options,
		})
	})
}

// Invocation is a function whose signature looks like:
//
//		func StartServer(server *http.Server) error {
//			return server.ListenAndServe()
//		}
//
// Like a constructor invocation may have unlimited count of arguments and
// they will be resolved automatically. The invocation can return an optional error.
// Error will be returned as is.
type Invocation interface{}

// Options group together container options.
//
//   account := di.Options(
//     di.Provide(NewAccountController),
//     di.Provide(NewAccountRepository),
//   )
//   auth := di.Options(
//     di.Provide(NewAuthController),
//     di.Provide(NewAuthRepository),
//   )
//   container, err := di.New(
//     account,
//     auth,
//   )
//   if err != nil {
//     // handle error
//   }
func Options(options ...Option) Option {
	return option(func(container *Container) {
		for _, opt := range options {
			opt.apply(container)
		}
	})
}

// WithLogger sets container logger.
func WithLogger(logger Logger) Option {
	return option(func(c *Container) {
		c.logger = logger
	})
}

// WithCompile ejects compile stage.
// Deprecated: Compile deprecated: https://github.com/goava/di/pull/9
func WithCompile() Option {
	return option(func(c *Container) {})
}

// CompileOption modifies compile behaviour.
// Deprecated: Compile deprecated: https://github.com/goava/di/pull/9
type CompileOption interface {
	apply(c *Container)
}

// ProvideParams is a Provide() method options. Name is a unique identifier of type instance. Provider is a constructor
// function. Interfaces is a interface that implements a provider result type.
type ProvideParams struct {
	Name        string
	Interfaces  []Interface
	IsPrototype bool
}

func (p ProvideParams) apply(params *ProvideParams) {
	*params = p
}

// InvokeOption is a functional option interface that modify invoke behaviour.
type InvokeOption interface {
	apply(params *InvokeParams)
}

// InvokeParams is a invoke parameters.
type InvokeParams struct {
	// The function
	Fn interface{}
}

func (p InvokeParams) apply(params *InvokeParams) {
	*params = p
}

// ResolveOption is a functional option interface that modify resolve behaviour.
type ResolveOption interface {
	apply(params *ResolveParams)
}

// Name specifies provider string identity. It needed when you have more than one
// definition of same type. You can identity type by name.
func Name(name string) ResolveOption {
	return resolveOption(func(params *ResolveParams) {
		params.Name = name
	})
}

// ResolveParams is a resolve parameters.
type ResolveParams struct {
	Name string
}

func (p ResolveParams) apply(params *ResolveParams) {
	*params = p
}

type option func(c *Container)

func (o option) apply(c *Container) { o(c) }

type provideOption func(params *ProvideParams)

func (o provideOption) apply(params *ProvideParams) {
	o(params)
}

type resolveOption func(params *ResolveParams)

func (o resolveOption) apply(params *ResolveParams) {
	o(params)
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
