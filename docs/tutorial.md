## Tutorial

Learn how to use the `di` package by building a simple application that processes HTTP requests.

The full tutorial code is available
[here](./../_examples/tutorial/main.go).

- [Tracing](#tracing)
- [Provide](#provide)
- [Resolve](#resolve)
- [Invoke](#invoke)
- [Lazy-loading](#lazy-loading)
- [Interfaces](#interfaces)
- [Groups](#groups)

### Tracing

Before starting, you can enable tracing to get more information about
the library lifecycle. The `di` package includes the default tracer that
prints output using the standard `log` package:

```go
func main() {
di.SetTracer(&di.StdTracer{})
//...
}
```

### Provide

First, we need to provide ways to build two fundamental
types: `http.Server` and `http.ServeMux`. Let's create simple
functional constructors that build them:

```go
// NewServer builds an HTTP server with the provided mux as handler.
func NewServer(mux *http.ServeMux) *http.Server {
return &http.Server{
Handler: mux,
}
}

// NewServeMux creates a new HTTP serve mux.
func NewServeMux() *http.ServeMux {
return &http.ServeMux{}
}
```

> Supported constructor signature:
>
> ```go
> // cleanup and error are optional
> func([dep1, dep2, depN]) (result, [cleanup, error])
> ```

Now we can teach the container to build these types in three ways:

Using the preferred functional option style:

```go
// create container
container, err := di.New(
    di.Provide(NewServer),
    di.Provide(NewServeMux),
)
if err != nil {
    // handle error
}
```

### Resolve

Next, we can resolve the built server from the container. To do this, define
the variable of the resolved type and pass the variable pointer to the `Resolve`
function.

If no error occurs, we can use the variable.

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

> The container creates singletons for combinations of the same type and
> tags.

### Invoke

As an alternative to resolve, we can use the `Invoke()` function of
the `Container`. It builds dependencies and calls the provided function. The Invoke
function can return an optional error.

```go
// StartServer starts the server.
func StartServer(server *http.Server) error {
return server.ListenAndServe()
}

if err := container.Invoke(StartServer); err != nil {
// handle error
}
```

Also, you can use the `di.Invoke()` container option to call some
initialization code.

```go
container, err := di.New(
    di.Provide(NewServer),
    di.Invoke(StartServer),
)
if err != nil {
// handle error
}
```

The container runs all `invoke functions` in the order they were
declared. If one of them fails, the compilation fails.

### Lazy-loading

Resulting dependencies will be lazy-loaded. If no one requests a type from
the container, it won't be constructed.

### Interfaces

You can provide an implementation as an interface. Use `di.As()` for this.
The arguments of this option must be a pointer(s) to an interface like
`new(http.Handler)`.

```go
di.Provide(NewServeMux, di.As(new(http.Handler)))
```

> This syntax with `new` can look strange, but I haven't found a better
> way to specify the interface.
>
> Create an issue if you know a better way ;)

Updated server constructor:

```go
// NewServer creates an HTTP server with the provided mux as handler.
func NewServer(handler http.Handler) *http.Server {
    return &http.Server{
        Handler: handler,
    }
}
```

Final code:

```go
container, err := di.New(
    // provide HTTP server
    di.Provide(NewServer),
    // provide HTTP serve mux as http.Handler interface
    di.Provide(NewServeMux, di.As(new(http.Handler)))
)
if err != nil {
    // handle error
}
```

Now the container uses `*http.ServeMux` as the implementation of `http.Handler`.
Interface usage contributes to writing more testable code.

### Groups

##### Grouping

The container automatically groups the same types into a `[]<type>` slice. It
works with `di.As()` too. For example, `di.As(new(http.Handler)`
automatically creates a group `[]http.Handler`.

Let's add some HTTP controllers using this feature. The main function of
controllers is registering routes. First, create an interface
for it.

```go
// Controller is an interface that can register its routes.
type Controller interface {
    RegisterRoutes(mux *http.ServeMux)
}
```

Next, create implementations for this interface.

##### Order implementation

```go
// OrderController is an HTTP controller for orders.
type OrderController struct {}

// NewOrderController creates an auth HTTP controller.
func NewOrderController() *OrderController {
    return &OrderController{}
}

// RegisterRoutes is a Controller interface implementation.
func (a *OrderController) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/orders", a.RetrieveOrders)
}

// RetrieveOrders loads orders and writes them to the writer.
func (a *OrderController) RetrieveOrders(writer http.ResponseWriter, request *http.Request) {
    // implementation
}
```

##### User implementation

```go
// UserController is an HTTP endpoint for users.
type UserController struct {}

// NewUserController creates a user HTTP endpoint.
func NewUserController() *UserController {
    return &UserController{}
}

// RegisterRoutes is a Controller interface implementation.
func (e *UserController) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/users", e.RetrieveUsers)
}

// RetrieveUsers loads users and writes them using the writer.
func (e *UserController) RetrieveUsers(writer http.ResponseWriter, request *http.Request) {
    // implementation
}
```

##### Container initialization code

Just like in the example with interfaces, we will use the `di.As()` provide
option.

```go
container, err := di.New(
    di.Provide(NewServer),        // provide HTTP server
    di.Provide(NewServeMux),       // provide HTTP serve mux
    // endpoints
    di.Provide(NewOrderController, di.As(new(Controller))),  // provide order controller
    di.Provide(NewUserController, di.As(new(Controller))),  // provide user controller
)
if err != nil {
    // handle error
}
```

Now we can use the `[]Controller` group in our mux. Updated code:

```go
// NewServeMux creates a new HTTP serve mux.
func NewServeMux(controllers []Controller) *http.ServeMux {
    mux := &http.ServeMux{}

    for _, controller := range controllers {
        controller.RegisterRoutes(mux)
    }

    return mux
}
```

The full tutorial code is available
[here](./../_examples/tutorial/main.go)