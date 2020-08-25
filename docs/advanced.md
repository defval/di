# Advanced features

- [Modules](#modules)
- [Named definitions](#named-definitions)
- [Optional parameters](#optional-parameters)
- [Struct fields injection](#struct-fields-injection)
- [Prototypes](#prototypes)
- [Cleanup](#cleanup)

### Modules

You can group previous options into single variable by using
`di.Options()` :

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

#### Specify name with `di.WithName()` *invoke option*:

```go
// provide leader database
di.Provide(NewLeader, di.WithName("leader"))
// provide follower database
di.Provide(NewFollower, di.WithName("follower"))
```

If you need to resolve it from the container use `di.Name()` *resolve
option*.

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

Also, `di.Inject` with tag `optional` provide ability to skip dependency
if it not exists in the container.

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

### Struct fields injection

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
d
### Prototypes

Use `di.Prototype()` option to create new instance for each resolve.

```go
di.Provide(NewRequestContext, di.Prototype())
```

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

