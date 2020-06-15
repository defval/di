DI Container
============
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff)](https://pkg.go.dev/github.com/goava/di)
[![Release](https://img.shields.io/github/tag/goava/di.svg?label=release&color=24B898&logo=github&style=for-the-badge)](https://github.com/goava/di/releases/latest)
[![Build Status](https://img.shields.io/travis/goava/di.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/goava/di)
[![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-green?style=for-the-badge)](https://goreportcard.com/report/github.com/goava/di)
[![Code Coverage](https://img.shields.io/codecov/c/github/goava/di.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/goava/di)

Dependency injection is one form of the broader technique of inversion
of control. It is used to increase modularity of the program and make it
extensible.

## Features

- Intuitive auto wiring based on go type system
- Interface implementations specification
- Same type groups for iteration
- Cleanup constructed instances
- Functional option interface
- Struct field injection
- Optional injection
- Type lazy-loading
- Named types

## Contents

- [Documentation](https://github.com/goava/di#documentation)
- [Install](https://github.com/goava/di#installing)
- [Tutorial](https://github.com/goava/di#tutorial)
  - [Provide](https://github.com/goava/di#provide)
  - [Resolve](https://github.com/goava/di#resolve)
  - [Invoke](https://github.com/goava/di#invoke)
  - [Lazy-loading](https://github.com/goava/di#lazy-loading)
  - [Interfaces](https://github.com/goava/di#interfaces)
  - [Groups](https://github.com/goava/di#groups)
- [Advanced features](https://github.com/goava/di#advanced-features)
  - [Modules](https://github.com/goava/di#modules)
  - [Named definitions](https://github.com/goava/di#named-definitions)
  - [Optional parameters](https://github.com/goava/di#optional-parameters)
  - [Struct field injection](https://github.com/goava/di#struct-fields-injection)
  - [Prototypes](https://github.com/goava/di#prototypes)
  - [Cleanup](https://github.com/goava/di#cleanup)
- [Comparison](https://github.com/goava/di#comparison)

## Documentation

You can use standard [pkg.go.dev](https://pkg.go.dev/github.com/goava/di) and inline code
comments or if you do not have experience with auto-wiring libraries
as [google/wire](https://github.com/google/wire),
[uber-go/dig](https://github.com/uber-go/dig) or another start with
[tutorial](https://github.com/goava/di#tutorial).

## Install

```shell
go get github.com/goava/di
```

## Tutorial

Let's learn to use `di` by example. We will code a simple application
that processes HTTP requests.

The full tutorial code is available [here](./_examples/tutorial/main.go).

### Provide

To start, we will need to provide way to build for two fundamental types: `http.Server`
and `http.ServeMux`. Let's create a simple functional constructors that build them:

```go
// NewServer builds a http server with provided mux as handler.
func NewServer(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Handler: mux,
	}
}

// NewServeMux creates a new http serve mux.
func NewServeMux() *http.ServeMux {
	return &http.ServeMux{}
}
```

> Supported constructor signature:
>
> ```go
> // cleanup and error is a optional
> func([dep1, dep2, depN]) (result, [cleanup, error])
> ```

Now we can teach the container to build these types in three ways:

In preferred functional option style:

```go
// create container
container, err := container.New(
	di.Provide(NewServer),
	di.Provide(NewServeMux),
)
if err != nil {
    // handle error
}
```

### Resolve

Next we can resolve the built server from the container. For this define the
variable of resolved type and pass variable pointer to `Resolve`
function.

If no error occurred we can use the variable.

```go
// declare type variable
var server *http.Server
// resolving
err := container.Resolve(&server)
if err != nil {
	// handle error
}

server.ListenAndServe()
```

> Note, by default the container creates singletons.
> But you can change this behaviour. See [Prototypes](https://github.com/goava/di#prototypes).

### Invoke

As an alternative to resolve we can use `Invoke()` function of `Container`. It builds
dependencies and calls provided function. Invoke function can return optional error.

```go
// StartServer starts the server.
func StartServer(server *http.Server) error {
    return server.ListenAndServe()
}

if err := container.Invoke(StartServer); err != nil {
	// handle error
}
```

Also you can use `di.Invoke()` container options for call some initialization code.

```go
container, err := di.New(
	di.Provide(NewServer),
	di.Invoke(StartServer),
)
if err != nil {
    // handle error
}
```

The container runs all `invoke functions` on the compile stage in the order they were declared. If one of then fails, the
compilation fails.

### Lazy-loading

Result dependencies will be lazy-loaded. If no one requires a type from
the container it won't be constructed.

### Interfaces

You can provide implementation as an interface. Use `di.As()` for it.
The arguments of this option must be a pointer(s) to an interface like `new(http.Handler)`.

```go
di.Provide(NewServeMux, di.As(new(http.Handler)))
```

> This syntax can look strange, but I haven't found a better way to
> specify the interface.

Updated server constructor:

```go
// NewServer creates a http server with provided mux as handler.
func NewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler: handler,
	}

```

Final code:

```go
container, err := di.New(
	// provide http server
	di.Provide(NewServer),
	// provide http serve mux as http.Handler interface
	di.Provide(NewServeMux, di.As(new(http.Handler)))
)
if err != nil {
    // handle error
}
```

Now container use `*http.ServeMux` as implementation of `http.Handler`. 
Interaface usage contributes to write more testable code.

### Groups

##### Grouping

Container automatically groups the same types to `[]<type>` slice. It works with `di.As()` too.
For example, `di.As(new(http.Handler)` automatically creates a group
`[]http.Handler`.

Let's add some http controllers using this feature. The main function of controllers is registering routes. At first, will create an interface for it.

```go
// Controller is an interface that can register its routes.
type Controller interface {
	RegisterRoutes(mux *http.ServeMux)
}
```

Next step is make implementation for this interface.

##### Order implementation

```go
// OrderController is a http controller for orders.
type OrderController struct {}

// NewOrderController creates a auth http controller.
func NewOrderController() *OrderController {
	return &OrderController{}
}

// RegisterRoutes is a Controller interface implementation.
func (a *OrderController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/orders", a.RetrieveOrders)
}

// RetrieveOrders loads orders and writes it to the writer.
func (a *OrderController) RetrieveOrders(writer http.ResponseWriter, request *http.Request) {
	// implementation
}
```

##### User implementation

```go
// UserController is a http endpoint for a user.
type UserController struct {}

// NewUserController creates a user http endpoint.
func NewUserController() *UserController {
	return &UserController{}
}

// RegisterRoutes is a Controller interface implementation.
func (e *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", e.RetrieveUsers)
}

// RetrieveUsers loads users and writes it using the writer.
func (e *UserController) RetrieveUsers(writer http.ResponseWriter, request *http.Request) {
    // implementation
}
```

##### Container initialization code

Just like in the example with interfaces, we will use `di.As()`
provide option.

```go
container, err := di.New(
	di.Provide(NewServer),        // provide http server
	di.Provide(NewServeMux),       // provide http serve mux
	// endpoints
	di.Provide(NewOrderController, di.As(new(Controller))),  // provide order controller
	di.Provide(NewUserController, di.As(new(Controller))),  // provide user controller
)
if err != nil {
    // handle error
}
```

Now we can use `[]Controller` group in our mux. Updated code:

```go
// NewServeMux creates a new http serve mux.
func NewServeMux(controllers []Controller) *http.ServeMux {
	mux := &http.ServeMux{}

	for _, controller := range controllers {
		controller.RegisterRoutes(mux)
	}

	return mux
}
```

The full tutorial code is available [here](./_examples/tutorial/main.go)

## Advanced features

### Modules

You can group previous options into single variable by using `di.Options()` :

```go
// account module
account := di.Options(
    di.Provide(NewAccountController), 
    di.Provide(NewAccountRepository),
)
// auth module
auth := di.Options(
    di.Provide(NewAuthController), 
    di.Provide(NewAuthRepository),
)
// build container
container, err := di.New(
    account, 
    auth,
)
if err != nil {
 // handle error
}
```

### Named definitions

If you have more than one instances of same type, you can specify alias. For example
two instances of database: leader - for writing, follower - for reading.

#### Wrap type into another unique type

```go
// Leader provides write database access.
type Leader struct {
	*Database
}

// Follower provides read database access.
type Follower struct {
	*Database
}
```

#### Specify name with `di.WithName()` *invoke option*:

```go
// provide leader database
di.Provide(NewLeader, di.WithName("leader"))
// provide follower database
di.Provide(NewFollower, di.WithName("follower"))
```

If you need to resolve it from the container use `di.Name()` *resolve option*.

```go
var db *Database
container.Resolve(&db, di.Name("leader"))
```

If you need to provide named definition in another constructor embed
`di.Inject`.

```go
// Parameters
type Parameters struct {
	di.Inject
	
	// use `di` tag for the container to know that field need to be injected.
	Leader *Database `di:"leader"`
	Follower *Database  `di:"follower"`
}

// NewService creates new service with provided parameters.
func NewService(parameters Parameters) *Service {
	return &Service{
		Leader:  parameters.Leader,
		Follower: parameters.Leader,
	}
}
```

### Optional parameters

Also, `di.Inject` with tag `optional` provide ability to skip dependency if it not exists
in the container.

```go
// ServiceParameter
type ServiceParameter struct {
	di.Inject
	
	Logger *Logger `di:"" optional:"true"`
}
```

> Constructors that declare dependencies as optional must handle the
> case of those dependencies being absent.

You can use naming and optional together.

```go
// ServiceParameter
type ServiceParameter struct {
	di.Inject
	
	StdOutLogger *Logger `di:"stdout"`
	FileLogger   *Logger `di:"file" optional:"true"`
}
```

### Struct field injection

To avoid constant constructor changes, you can use `di.Inject`. Only
struct pointers are supported as constructing result. And only 
`di`-taged fields will be injected. Such a constructor will work with 
using `di` only.

```go
// Controller has some endpoints.
type Controller struct {
    di.Inject // enables struct field injection 

    // fields must be public and have tag di
    // tag lets to specify fields need to be injected
    Users   UserService     `di:""`
    Friends FriendsService  `di:""`
}

// NewController creates controller.
func NewController() *Controller {
    return &Controller{}
}
```

### Prototypes

Use `di.Prototype()` option to create new instance for each resolve.

```go
di.Provide(NewRequestContext, di.Prototype())
```

### Cleanup

If the constructor creates a value that needs to be cleaned up, then it can
return a closure to clean up the resource.

```go
func NewFile(log Logger, path Path) (*os.File, func(), error) {
    f, err := os.Open(string(path))
    if err != nil {
        return nil, nil, err
    }
    cleanup := func() {
        if err := f.Close(); err != nil {
            log.Log(err)
        }
    }
    return f, cleanup, nil
}
```

After `container.Cleanup()` call, it iterates over instances and calls
cleanup function if it exists.

```go
container, err := di.New(
	// ...
    di.Provide(NewFile),
)
if err != nil {
    // handle error
}
// do something
container.Cleanup() // file was closed
```