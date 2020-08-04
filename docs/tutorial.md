## Tutorial

Let's learn to use `di` by example. We will code a simple application
that processes HTTP requests.

The full tutorial code is available [here](./../_examples/tutorial/main.go).

- [Provide](#provide)
- [Resolve](#resolve)
- [Invoke](#invoke)
- [Lazy-loading](#lazy-loading)
- [Interfaces](#interfaces)
- [Groups](#groups)

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
> But you can change this behaviour. See [Prototypes](.#prototypes).

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

The full tutorial code is available [here](./../_examples/tutorial/main.go)