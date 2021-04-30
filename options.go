package di

// Option is a functional option that configures container. If you don't know about functional
// options, see https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis.
// Below presented all possible options with their description:
//
// 	- di.Provide - provide constructors
//	- di.Invoke - add invocations
//	- di.Resolve - resolves type
type Option interface {
	apply(c *diopts)
}

// Provide returns container option that provides to container reliable way to build type. The constructor will
// be invoked lazily on-demand. For more information about constructors see Constructor interface. ProvideOption can
// add additional behavior to the process of type resolving.
func Provide(constructor Constructor, options ...ProvideOption) Option {
	frame := stacktrace(0)
	return option(func(c *diopts) {
		c.provides = append(c.provides, provideOptions{
			frame,
			constructor,
			options,
		})
	})
}

// ProvideValue provides value as is.
func ProvideValue(value Value, options ...ProvideOption) Option {
	frame := stacktrace(0)
	return option(func(c *diopts) {
		c.values = append(c.values, provideValueOptions{
			frame,
			value,
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

// Value
type Value interface{}

// ProvideOption is a functional option interface that modify provide behaviour. See di.As(), di.WithName().
type ProvideOption interface {
	applyProvide(params *ProvideParams)
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
// Deprecated: use di.Tags.
func WithName(name string) ProvideOption {
	return provideOption(func(params *ProvideParams) {
		if params.Tags == nil {
			params.Tags = Tags{}
		}
		params.Tags["name"] = name
	})
}

// Decorator can modify container instance.
// EXPERIMENTAL FEATURE: functional can be changed.
type Decorator func(value Value) error

// Decorate will be called after type construction. You can modify your pointer types.
// EXPERIMENTAL FEATURE: functional can changed.
func Decorate(decorators ...Decorator) ProvideOption {
	return provideOption(func(params *ProvideParams) {
		params.Decorators = append(params.Decorators, decorators...)
	})
}

// Resolve returns container options that resolves type into target. All resolves will be done on compile stage
// after call invokes.
func Resolve(target Pointer, options ...ResolveOption) Option {
	frame := stacktrace(0)
	return option(func(c *diopts) {
		c.resolves = append(c.resolves, resolveOptions{
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
	frame := stacktrace(0)
	return option(func(c *diopts) {
		c.invokes = append(c.invokes, invokeOptions{
			frame,
			fn,
			options,
		})
	})
}

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
	return option(func(container *diopts) {
		for _, opt := range options {
			opt.apply(container)
		}
	})
}

// ProvideParams is a Provide() method options. Name is a unique identifier of type instance. Provider is a constructor
// function. Interfaces is a interface that implements a provider result type.
type ProvideParams struct {
	Tags       Tags
	Interfaces []Interface
	Decorators []Decorator
}

func (p ProvideParams) applyProvide(params *ProvideParams) {
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
	applyResolve(params *ResolveParams)
}

// Name specifies provider string identity. It needed when you have more than one
// definition of same type. You can identity type by name.
// Deprecated: use di.Tags
func Name(name string) ResolveOption {
	return resolveOption(func(params *ResolveParams) {
		if params.Tags == nil {
			params.Tags = Tags{}
		}
		params.Tags["name"] = name
	})
}

// ResolveParams is a resolve parameters.
type ResolveParams struct {
	Tags Tags
}

func (p ResolveParams) applyResolve(params *ResolveParams) {
	*params = p
}

type option func(c *diopts)

func (o option) apply(c *diopts) { o(c) }

type provideOption func(params *ProvideParams)

func (o provideOption) applyProvide(params *ProvideParams) {
	o(params)
}

type resolveOption func(params *ResolveParams)

func (o resolveOption) applyResolve(params *ResolveParams) {
	o(params)
}

// struct that contains constructor with options.
type provideOptions struct {
	frame       callerFrame
	constructor Constructor
	options     []ProvideOption
}

// struct that contains value with options.
type provideValueOptions struct {
	frame   callerFrame
	value   Value
	options []ProvideOption
}

// struct that contains invoke function with options.
type invokeOptions struct {
	frame   callerFrame
	fn      Invocation
	options []InvokeOption
}

// struct that container resolve target with options.
type resolveOptions struct {
	frame   callerFrame
	target  Pointer
	options []ResolveOption
}
