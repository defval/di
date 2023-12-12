# Advanced Features

- [Modules](#modules)
- [Tags](#tags)
- [ProvideValue](#providevalue)
- [Optional Parameters](#optional-parameters)
- [Struct Field Injection](#struct-field-injection)
- [Iteration](#iteration)
- [Decoration](#decoration)
- [Cleanup](#cleanup)
- [Container Chaining / Scopes](#container-chaining--scopes)

### Modules

You can group previous options into a single variable using `di.Options()`:

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

### Tags

If you have more than one instance of the same type, you can specify an alias.
For example, two instances of a database: leader - for writing, follower -
for reading.

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

#### Specify tags with `di.Tags` *provide option*:

```go
// provide leader database
di.Provide(NewLeader, di.Tags{"type":"leader"})
// provide follower database
di.Provide(NewFollower, di.Tags{"type", "follower"}))
```

If you need to resolve it from the container, use `di.Tags` *resolve
option*.

```go
var db *Database
container.Resolve(&db, di.Tags{"type": "leader"}))
```

If you need to provide a named definition in another constructor, embed
`di.Inject`.

```go
// Parameters
type Parameters struct {
	di.Inject
	
	// use tag for the container to know that field need to be injected.
	Leader *Database `di:"type=leader"`
	Follower *Database  `di:"type=follower"`
}

// NewService creates a new service with provided parameters.
func NewService(parameters Parameters) *Service {
	return &Service{
		Leader:  parameters.Leader,
		Follower: parameters.Leader,
	}
}
```

If you need to resolve all types with the same tag key, use `*` as the tag
value:

```go
var db []*Database
di.Resolve(&db, di.Tags{"type": "*"})
```

### ProvideValue

Instead of using `di.Provide` to provide a constructor, you can use `di.ProvideValue` and provide values directly.
This is useful to provide primitive values or values that are easily constructed. You can combine this with the use
of `di.Tags` if you have multiple values of the same type and want to identify each one.

```go
di.New(
    di.ProvideValue(time.Duration(10*time.Second), di.Tags{"name": "http-timeout"}),
)

var timeout time.Duration
c.Resolve(&timeout, di.Tags{"name": "http-timeout"})
```

To differentiate between multiple values of the same type, you can also use golang type aliases instead of using
`di.Tags`.

```go
type ProjectName string
type ProjectVersion string

c, err := di.New(
    di.ProvideValue(ProjectName("my-project")),
    di.ProvideValue(ProjectVersion("1.0.0")),
)

var pn ProjectName
c.Resolve(&pn)
var pv ProjectVersion
c.Resolve(&pv)
```

### Optional Parameters

Also, `di.Inject` with tag `di:"optional"` provides the ability to skip a dependency
if it does not exist in the container.

```go
// ServiceParameter
type ServiceParameter struct {
	di.Inject
	
	Logger *Logger `di:"optional"`
}
```

> Constructors that declare dependencies as optional must handle the
> case of those dependencies being absent.

You can use tagged and optional together.

```go
// ServiceParameter
type ServiceParameter struct {
	di.Inject
	
	StdOutLogger *Logger `di:"type=stdout"`
	FileLogger   *Logger `di:"type=file,optional"`
}
```

If you need to skip field injection, use `di:"skip"` tags for this:

```go
// ServiceParameter
type ServiceParameter struct {
	di.Inject
	
	StdOutLogger *Logger    `di:"type=stdout"`
	FileLogger   *Logger    `di:"type=file,optional"`
	SkipField    *SomeType  `di:"skip"` // injection skipped
}
```

### Struct Field Injection

To avoid constant constructor changes, you can use `di.Inject`. Only
struct pointers are supported as constructing results. And only
`di`-tagged fields will be injected. Such a constructor will work with
using `di` only.

```go
// Controller has some endpoints.
type Controller struct {
    di.Inject // enables struct field injection 

    // fields must be public
    // tag lets to specify fields need to be injected
    Users   UserService
    Friends FriendsService  `di:"type=cached"`
}

// NewController creates a controller.
func NewController() *Controller {
    return &Controller{}
}
```
### Iteration

The `di` package provides iteration capabilities, allowing you to iterate over a group of a specific Pointer type with the `IterateFunc`. This can be useful when working with multiple instances of a type or when you need to perform actions on each instance.

```go
// ValueFunc is a lazy-loading wrapper for iteration.
type ValueFunc func() (interface{}, error)

// IterateFunc is a function that will be called on each instance in the iterate selection.
type IterateFunc func(tags Tags, value ValueFunc) error
```

To use iteration with the container, follow the example below:

```go
var servers []*http.Server
iterFn := func(tags di.Tags, loader ValueFunc) error {
	i, err := loader()
	if err != nil {
		return err
	}
	// do stuff with result: i.(*http.Server)
	return nil
}

container.Iterate(&servers, iterFn)
```

In this example, the `Iterate` method is called on the container, passing a slice of pointers to the desired type (in this case, `*http.Server`) and the iterate function, which will be executed on each instance.

### Decoration

The `di` package supports decoration, allowing you to modify container instances through the use of decorators. This can be helpful when you need to make additional modifications to instances after they have been constructed.

```go
// Decorator can modify container instance.
type Decorator func(value Value) error

// Decorate will be called after type construction. You can modify your pointer types.
func Decorate(decorators ...Decorator) ProvideOption {
	return provideOption(func(params *ProvideParams) {
		params.Decorators = append(params.Decorators, decorators...)
	})
}
```

To use decorators, you can add them to the `Provide` method using the `Decorate` function. Here's an example of a decorator that logs the creation of each instance:

```go
// Logger is a simple logger interface for demonstration purposes
type Logger interface {
	Log(message string)
}

// logInstanceCreation is a decorator that logs the creation of instances
func logInstanceCreation(logger Logger) Decorator {
	return func(value Value) error {
		logger.Log(fmt.Sprintf("Instance of type logger created"))
		return nil
	}
}

// Usage example
container, err := di.New(
	di.Provide(NewMyType, di.Decorate(logInstanceCreation(myLogger))),
)
```

In this example, the `logInstanceCreation` decorator logs a message every time a new instance is created. The decorator is added to the `Provide` method using the `Decorate` function, and it is executed after the type construction.

### Cleanup

If the constructor creates a value that needs to be cleaned up, then it
can return a closure to clean up the resource.

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

### Container Chaining / Scopes

You can chain containers together so that values can be resolved from a
parent container. This lets you do things like have a configuration
scope container and an application scoped container. By keeping
configuration values in a different container, you can re-create
the application scoped container when you make configuration changes
since each container has an independent lifecycle.

**Note:** You should cleanup each container manually.

```go
configContainer, err := container.New(
    di.Provide(NewServerConfig),
)

appContainer, err := container.New(di.Provide(config *SeverConfig) *http.Server {
    sever := ...
    return server
})

if err := appContainer.AddParent(configContainer); err != nil {
   // handle error
}

var server *http.Server
err := appContainer.Resolve(&server)
```