# Advanced features

- [Modules](#modules)
- [Tags](#tags)
- [Optional parameters](#optional-parameters)
- [Struct fields injection](#struct-fields-injection)
- [Iteration](#iteration)
- [Cleanup](#cleanup)
- [Container Chaining / Scopes](#container-chaining--scopes)

### Modules

You can group previous options into single variable by using
`di.Options()`:

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

If you have more than one instances of same type, you can specify alias.
For example two instances of database: leader - for writing, follower -
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

If you need to resolve it from the container use `di.Tags` *resolve
option*.

```go
var db *Database
container.Resolve(&db, di.Tags{"type": "leader"}))
```

If you need to provide named definition in another constructor embed
`di.Inject`.

```go
// Parameters
type Parameters struct {
	di.Inject
	
	// use tag for the container to know that field need to be injected.
	Leader *Database `di:"type=leader"`
	Follower *Database  `di:"type=follower"`
}

// NewService creates new service with provided parameters.
func NewService(parameters Parameters) *Service {
	return &Service{
		Leader:  parameters.Leader,
		Follower: parameters.Leader,
	}
}
```

If you need to resolve all types with same tag key, use `*` as tag
value:

```go
var db []*Database
di.Resolve(&db, di.Tags{"type": "*"})
```

### Optional parameters

Also, `di.Inject` with tag `di:"optional"` provide ability to skip dependency
if it not exists in the container.

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

If you need to skip fields injection use `di:"skip"` tags for this:

```go
// ServiceParameter
type ServiceParameter struct {
	di.Inject
	
	StdOutLogger *Logger    `di:"type=stdout"`
	FileLogger   *Logger    `di:"type=file,optional"`
	SkipField    *SomeType  `di:"skip"` // injection skipped
}
```

### Struct fields injection

To avoid constant constructor changes, you can use `di.Inject`. Only
struct pointers are supported as constructing result. And only
`di`-taged fields will be injected. Such a constructor will work with
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

// NewController creates controller.
func NewController() *Controller {
    return &Controller{}
}
```

### Iteration

TBD

### Decoration

TBD

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
scope container and an application scoped container.  By keeping 
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

