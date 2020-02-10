package di_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/goava/di"
	"github.com/goava/di/internal/ditest"
)

func TestContainerCompileErrors(t *testing.T) {
	t.Run("dependency cycle cause panic", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewCycleFooBar)
		c.MustProvide(ditest.NewBar)
		c.MustCompileError("the graph cannot be cyclic")
	})

	t.Run("not existing dependency cause compile error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewBar)
		c.MustCompileError("*ditest.Bar: dependency *ditest.Foo not exists in container")
	})

	t.Run("not existing non pointer dependency cause compile error", func(t *testing.T) {
		c := NewTestContainer(t)
		type TestStruct struct {
		}

		c.MustProvide(func(s TestStruct) bool {
			return true
		})

		require.PanicsWithValue(t, "bool: dependency di_test.TestStruct not exists in container", func() {
			c.Compile()
		})
	})
}

func TestContainerProvideErrors(t *testing.T) {
	t.Run("provide string cause panic", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvideError("string", "The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `string`")
	})

	t.Run("provide nil cause panic", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvideError(nil, "The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `nil`")
	})

	t.Run("provide struct pointer cause panic", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvideError(&ditest.Foo{}, "The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `*ditest.Foo`")
	})

	t.Run("provide constructor without result cause panic", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvideError(ditest.ConstructorWithoutResult, "The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `github.com/goava/di/internal/ditest.ConstructorWithoutResult`")
	})

	t.Run("provide constructor with many results cause panic", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvideError(ditest.ConstructorWithManyResults, "The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `github.com/goava/di/internal/ditest.ConstructorWithManyResults`")
	})

	t.Run("provide constructor with incorrect result error argument", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvideError(ditest.ConstructorWithIncorrectResultError, "The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `github.com/goava/di/internal/ditest.ConstructorWithIncorrectResultError`")
	})

	t.Run("provide duplicate", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustProvideError(ditest.NewFoo, "The `*ditest.Foo` type already exists in container")
	})

	t.Run("provide as not implemented interface cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustProvideError(ditest.NewBar, "*ditest.Bar not implement ditest.Barer", new(ditest.Barer))
	})

	t.Run("provide as not interface cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustProvideError(ditest.NewBar, "*ditest.Foo: not a pointer to interface", new(ditest.Foo))
	})
}

func TestContainerExtractErrors(t *testing.T) {
	t.Run("container panic on trying extract before compilation", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := &ditest.Foo{}
		c.MustProvide(ditest.CreateFooConstructor(foo))
		var extracted *ditest.Foo
		c.MustExtractError(&extracted, "container not compiled")
	})

	t.Run("extract into string cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustCompile()
		c.MustExtractError("string", "resolve target must be a pointer, got `string`")
	})

	t.Run("extract into struct cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustCompile()
		c.MustExtractError(struct{}{}, "resolve target must be a pointer, got `struct {}`")
	})

	t.Run("extract into nil cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustCompile()
		c.MustExtractError(nil, "resolve target must be a pointer, got `nil`")
	})

	t.Run("container does not find type because its named", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := &ditest.Foo{}
		c.MustProvideWithName("foo", ditest.CreateFooConstructor(foo))
		c.MustCompile()

		var extracted *ditest.Foo
		c.MustExtractError(&extracted, "*ditest.Foo: not exists in container")
	})

	t.Run("extract returns error because dependency constructing failed", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.CreateFooConstructorWithError(errors.New("internal error")))
		c.MustProvide(ditest.NewBar)
		c.MustCompile()
		var bar *ditest.Bar
		c.MustExtractError(&bar, "*ditest.Foo: internal error")
	})

	t.Run("extract interface with multiple implementations cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustProvide(ditest.NewBar, new(ditest.Fooer))
		c.MustProvide(ditest.NewBaz, new(ditest.Fooer))
		c.MustCompile()

		var extracted ditest.Fooer
		c.MustExtractError(&extracted, "ditest.Fooer: have several implementations")
	})
}

func TestContainerInvokeErrors(t *testing.T) {
	t.Run("addInvocation function with incorrect signature cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustCompile()
		c.MustInvokeError(func() *ditest.Foo {
			return nil
		}, "the addInvocation function must be a function like `func([dep1, dep2, ...]) [error]`, got `func() *ditest.Foo`")
	})

	t.Run("addInvocation function with undefined dependency cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustCompile()
		c.MustInvokeError(func(foo *ditest.Foo) {}, "resolve invocation (github.com/goava/di_test.TestContainerInvokeErrors.func2.1): *ditest.Foo: not exists in container")
	})

	t.Run("addInvocation before compile cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustInvokeError(func() {}, "container not compiled")
	})
}

func TestContainerProvide(t *testing.T) {
	t.Run("container successfully accept simple constructor", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
	})

	t.Run("container successfully accept constructor with error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.CreateFooConstructorWithError(nil))
	})

	t.Run("container successfully accept constructor with cleanup function", func(t *testing.T) {
		c := NewTestContainer(t)

		cleanup := func() {}
		c.MustProvide(ditest.CreateFooConstructorWithCleanup(cleanup))
	})

}

func TestContainerExtract(t *testing.T) {
	t.Run("container extract correct pointer", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := &ditest.Foo{}
		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustCompile()

		var extracted *ditest.Foo
		c.MustExtractPtr(foo, &extracted)
	})

	t.Run("container extract same pointer on each extraction", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := &ditest.Foo{}
		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustCompile()

		var extracted1 *ditest.Foo
		c.MustExtractPtr(foo, &extracted1)

		var extracted2 *ditest.Foo
		c.MustExtractPtr(foo, &extracted2)
	})

	t.Run("container extract instance if error is nil", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.CreateFooConstructorWithError(nil))
		c.MustCompile()

		var extracted *ditest.Foo
		c.MustExtract(&extracted)
	})

	t.Run("container extract instance if cleanup and error is nil", func(t *testing.T) {
		c := NewTestContainer(t)

		c.MustProvide(ditest.CreateFooConstructorWithCleanupAndError(nil, nil))
		c.MustCompile()

		var extracted *ditest.Foo
		c.MustExtract(&extracted)
	})

	t.Run("container extract correct named pointer", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := &ditest.Foo{}
		c.MustProvideWithName("foo", ditest.CreateFooConstructor(foo))
		c.MustCompile()

		var extracted *ditest.Foo
		c.MustExtractWithName("foo", &extracted)
	})

	t.Run("container extract correct interface implementation", func(t *testing.T) {
		c := NewTestContainer(t)
		bar := &ditest.Bar{}
		c.MustProvide(ditest.NewFoo)
		c.MustProvide(ditest.CreateBarConstructor(bar), new(ditest.Fooer))
		c.MustCompile()

		var extracted ditest.Fooer
		c.MustExtractPtr(bar, &extracted)
	})

	t.Run("container creates group from interface and extract it", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustProvide(ditest.NewBar, new(ditest.Fooer))
		c.MustProvide(ditest.NewBaz, new(ditest.Fooer))
		c.MustCompile()

		var group []ditest.Fooer
		c.MustExtract(&group)
		require.Len(t, group, 2)
	})

	t.Run("container extract new instance of prototype by each extraction", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustProvidePrototype(ditest.NewBar)
		c.MustCompile()

		var extracted1 *ditest.Bar
		c.MustExtract(&extracted1)
		var extracted2 *ditest.Bar
		c.MustExtract(&extracted2)

		c.MustNotEqualPointer(extracted1, extracted2)
	})
}

func TestContainerResolve(t *testing.T) {
	t.Run("container resolve correct argument", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := &ditest.Foo{}
		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustProvide(ditest.NewBar)
		c.MustCompile()

		var bar *ditest.Bar
		c.MustExtract(&bar)
		c.MustEqualPointer(foo, bar.Foo())
	})

	t.Run("container resolve correct interface implementation", func(t *testing.T) {
		c := NewTestContainer(t)

		foo := ditest.NewFoo()
		bar := ditest.NewBar(foo)

		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustProvide(ditest.CreateBarConstructor(bar), new(ditest.Fooer))
		c.MustProvide(ditest.NewQux)
		c.MustCompile()

		var qux *ditest.Qux
		c.MustExtract(&qux)
		c.MustEqualPointer(bar, qux.Fooer())
	})

	t.Run("container resolve correct group", func(t *testing.T) {
		c := NewTestContainer(t)

		c.MustProvide(ditest.NewFoo)
		c.MustProvide(ditest.NewBar, new(ditest.Fooer))
		c.MustProvide(ditest.NewBaz, new(ditest.Fooer))
		c.MustProvide(ditest.NewFooerGroup)
		c.MustCompile()

		var bar *ditest.Bar
		c.MustExtract(&bar)

		var baz *ditest.Baz
		c.MustExtract(&baz)

		var group *ditest.FooerGroup
		c.MustExtract(&group)
		require.Len(t, group.Fooers(), 2)
		c.MustEqualPointer(bar, group.Fooers()[0])
		c.MustEqualPointer(baz, group.Fooers()[1])
	})
}

func TestContainerResolveEmbedParameters(t *testing.T) {
	t.Run("container resolve embed parameters", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := ditest.NewFoo()
		bar := ditest.NewBar(foo)
		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustProvide(ditest.CreateBarConstructor(bar))
		c.MustProvide(ditest.NewBazFromParameters)
		c.MustCompile()

		var extracted *ditest.Baz
		c.MustExtract(&extracted)
		c.MustEqualPointer(foo, extracted.Foo())
		c.MustEqualPointer(bar, extracted.Bar())
	})

	t.Run("container skip optional parameter", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := ditest.NewFoo()
		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustProvide(ditest.NewBazFromParameters)
		c.MustCompile()

		var extracted *ditest.Baz
		c.MustExtract(&extracted)
		c.MustEqualPointer(foo, extracted.Foo())
		require.Nil(t, extracted.Bar())
	})

	t.Run("container resolve optional not existing group as nil", func(t *testing.T) {
		c := NewTestContainer(t)
		type Params struct {
			di.Parameter
			Handlers []http.Handler `di:"optional"`
		}
		c.MustProvide(func(params Params) bool {
			return params.Handlers == nil
		})
		c.MustCompile()
		var extracted bool
		c.MustExtract(&extracted)
		require.True(t, extracted)
	})

	t.Run("container skip private fields in parameter", func(t *testing.T) {
		c := NewTestContainer(t)
		type Param struct {
			di.Parameter
			private    []http.Handler `di:"optional"`
			Addrs      []net.Addr     `di:"optional"`
			HaveNotTag string
		}
		c.MustProvide(func(param Param) bool {
			return param.Addrs == nil
		})
		c.MustCompile()
		var extracted bool
		c.MustExtract(&extracted)
		require.True(t, extracted)
	})
}

func TestContainerInvoke(t *testing.T) {
	t.Run("container call addInvocation function", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustCompile()
		var invokeCalled bool
		c.MustInvoke(func() {
			invokeCalled = true
		})
		require.True(t, invokeCalled)
	})

	t.Run("container resolve dependencies in addInvocation function", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := ditest.NewFoo()
		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustCompile()
		c.MustInvoke(func(invokeFoo *ditest.Foo) {
			c.MustEqualPointer(foo, invokeFoo)
		})
	})

	t.Run("container addInvocation return correct error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.Compile()
		c.MustInvokeError(func(foo *ditest.Foo) error {
			return errors.New("addInvocation error")
		}, "addInvocation error")
	})

	t.Run("container addInvocation with nil error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.Compile()
		c.MustInvoke(func(foo *ditest.Foo) error {
			return nil
		})
	})
}

func TestContainerResolveParameterBag(t *testing.T) {
	t.Run("container extract correct parameter bag for type", func(t *testing.T) {
		c := NewTestContainer(t)

		c.Provide(ditest.NewFooWithParameters, di.ProvideParams{
			Parameters: di.ParameterBag{
				"name": "test",
			},
		})

		c.MustCompile()

		var foo *ditest.Foo
		err := c.Resolve(&foo)

		require.NoError(t, err)
		require.Equal(t, "test", foo.Name)
	})

	t.Run("container extract correct parameter bag for named type", func(t *testing.T) {
		c := NewTestContainer(t)

		c.Provide(ditest.NewFooWithParameters, di.ProvideParams{
			Name: "named",
			Parameters: di.ParameterBag{
				"name": "test",
			},
		})

		c.MustCompile()

		var foo *ditest.Foo
		err := c.Resolve(&foo, di.ExtractParams{
			Name: "named",
		})

		require.NoError(t, err)
		require.Equal(t, "test", foo.Name)
	})
}

func TestContainerCleanup(t *testing.T) {
	t.Run("cleanup container", func(t *testing.T) {
		c := NewTestContainer(t)
		var cleanupCalled bool
		c.MustProvide(ditest.CreateFooConstructorWithCleanup(func() { cleanupCalled = true }))
		c.MustCompile()

		var extracted *ditest.Foo
		c.MustExtract(&extracted)
		c.Cleanup()

		require.True(t, cleanupCalled)
	})

	t.Run("cleanup run in correct order", func(t *testing.T) {
		c := NewTestContainer(t)
		var cleanupCalls []string
		c.MustProvide(func(bar *ditest.Bar) (*ditest.Foo, func()) {
			return &ditest.Foo{}, func() { cleanupCalls = append(cleanupCalls, "foo") }
		})
		c.MustProvide(func() (*ditest.Bar, func()) {
			return &ditest.Bar{}, func() { cleanupCalls = append(cleanupCalls, "bar") }
		})
		c.MustCompile()

		var foo *ditest.Foo
		c.MustExtract(&foo)
		c.Cleanup()
		require.Equal(t, []string{"bar", "foo"}, cleanupCalls)
	})

	t.Run("cleanup for every prototyped instance", func(t *testing.T) {
		c := NewTestContainer(t)
		var cleanupCalls []string
		c.Provide(func() (*ditest.Foo, func()) {
			return &ditest.Foo{}, func() {
				cleanupCalls = append(cleanupCalls, fmt.Sprintf("foo_%d", len(cleanupCalls)))
			}
		}, di.ProvideParams{
			IsPrototype: true,
		})
		c.MustCompile()
		var foo1, foo2 *ditest.Foo
		c.MustExtract(&foo1)
		c.MustExtract(&foo2)
		c.Cleanup()
		require.Equal(t, []string{"foo_0", "foo_1"}, cleanupCalls)
	})
}

func TestContainer_GraphVisualizing(t *testing.T) {
	t.Run("graph", func(t *testing.T) {
		c := NewTestContainer(t)

		c.MustProvide(ditest.NewLogger)
		c.MustProvide(ditest.NewServer)
		c.MustProvide(ditest.NewRouter, new(http.Handler))
		c.MustProvide(ditest.NewAccountController, new(ditest.Controller))
		c.MustProvide(ditest.NewAuthController, new(ditest.Controller))
		c.MustCompile()

		var graph *di.Graph
		require.NoError(t, c.Resolve(&graph))

		fmt.Println(graph.String())

		require.Equal(t, `digraph  {
	subgraph cluster_s3 {
		ID = "cluster_s3";
		bgcolor="#E8E8E8";color="lightgrey";fontcolor="#46494C";fontname="COURIER";label="";style="rounded";
		n9[color="#46494C",fontcolor="white",fontname="COURIER",label="*di.Graph",shape="box",style="filled"];
		
	}subgraph cluster_s2 {
		ID = "cluster_s2";
		bgcolor="#E8E8E8";color="lightgrey";fontcolor="#46494C";fontname="COURIER";label="";style="rounded";
		n6[color="#46494C",fontcolor="white",fontname="COURIER",label="*ditest.AccountController",shape="box",style="filled"];
		n8[color="#46494C",fontcolor="white",fontname="COURIER",label="*ditest.AuthController",shape="box",style="filled"];
		n7[color="#E54B4B",fontcolor="white",fontname="COURIER",label="[]ditest.Controller",shape="doubleoctagon",style="filled"];
		n4[color="#E5984B",fontcolor="white",fontname="COURIER",label="ditest.RouterParams",shape="box",style="filled"];
		
	}subgraph cluster_s0 {
		ID = "cluster_s0";
		bgcolor="#E8E8E8";color="lightgrey";fontcolor="#46494C";fontname="COURIER";label="";style="rounded";
		n1[color="#46494C",fontcolor="white",fontname="COURIER",label="*log.Logger",shape="box",style="filled"];
		
	}subgraph cluster_s1 {
		ID = "cluster_s1";
		bgcolor="#E8E8E8";color="lightgrey";fontcolor="#46494C";fontname="COURIER";label="";style="rounded";
		n3[color="#46494C",fontcolor="white",fontname="COURIER",label="*http.ServeMux",shape="box",style="filled"];
		n2[color="#46494C",fontcolor="white",fontname="COURIER",label="*http.Server",shape="box",style="filled"];
		n5[color="#2589BD",fontcolor="white",fontname="COURIER",label="http.Handler",style="filled"];
		
	}splines="ortho";
	n6->n7[color="#949494"];
	n8->n7[color="#949494"];
	n3->n5[color="#949494"];
	n1->n2[color="#949494"];
	n1->n3[color="#949494"];
	n1->n6[color="#949494"];
	n1->n8[color="#949494"];
	n7->n4[color="#949494"];
	n4->n3[color="#949494"];
	n5->n2[color="#949494"];
	
}`, graph.String())
	})
}

// NewTestContainer
func NewTestContainer(t *testing.T) *TestContainer {
	return &TestContainer{t, di.New()}
}

// TestContainer
type TestContainer struct {
	t *testing.T
	*di.Container
}

func (c *TestContainer) MustProvide(provider interface{}, as ...di.Interface) {
	require.NotPanics(c.t, func() {
		c.Provide(provider, di.ProvideParams{
			Interfaces: as,
		})
	}, "provide should not panic")
}

func (c *TestContainer) MustProvidePrototype(provider interface{}, as ...di.Interface) {
	require.NotPanics(c.t, func() {
		c.Provide(provider, di.ProvideParams{
			Interfaces:  as,
			IsPrototype: true,
		})
	})
}

func (c *TestContainer) MustProvideWithName(name string, provider interface{}, as ...di.Interface) {
	require.NotPanics(c.t, func() {
		c.Provide(provider, di.ProvideParams{
			Name:       name,
			Interfaces: as,
		})
	})
}

func (c *TestContainer) MustProvideError(provider interface{}, msg string, as ...di.Interface) {
	require.PanicsWithValue(c.t, msg, func() {
		c.Provide(provider, di.ProvideParams{
			Interfaces: as,
		})
	})
}

func (c *TestContainer) MustCompile() {
	require.NotPanics(c.t, func() {
		c.Compile()
	})
}

func (c *TestContainer) MustCompileError(msg string) {
	require.PanicsWithValue(c.t, msg, func() {
		c.Compile()
	})
}

func (c *TestContainer) MustExtract(target interface{}) {
	require.NoError(c.t, c.Resolve(target))
}

func (c *TestContainer) MustExtractWithName(name string, target interface{}) {
	require.NoError(c.t, c.Resolve(target, di.ExtractParams{
		Name: name,
	}))
}

func (c *TestContainer) MustExtractError(target interface{}, msg string) {
	require.EqualError(c.t, c.Resolve(target, di.ExtractParams{}), msg)
}

func (c *TestContainer) MustExtractWithNameError(name string, target interface{}, msg string) {
	require.EqualError(c.t, c.Resolve(target, di.ExtractParams{
		Name: name,
	}), msg)
}

// MustExtractPtr extract value from container into target and check that target and expected pointers are equal.
func (c *TestContainer) MustExtractPtr(expected, target interface{}) {
	c.MustExtract(target)

	// indirect
	actual := reflect.ValueOf(target).Elem().Interface()
	c.MustEqualPointer(expected, actual)
}

func (c *TestContainer) MustExtractPtrWithName(expected interface{}, name string, target interface{}) {
	c.MustExtractWithName(name, target)

	actual := reflect.ValueOf(target).Elem().Interface()
	c.MustEqualPointer(expected, actual)
}

func (c *TestContainer) MustInvoke(fn interface{}) {
	require.NoError(c.t, c.Invoke(fn))
}

func (c *TestContainer) MustInvokeError(fn interface{}, msg string) {
	require.EqualError(c.t, c.Invoke(fn), msg)
}

func (c *TestContainer) MustEqualPointer(expected interface{}, actual interface{}) {
	require.Equal(c.t,
		fmt.Sprintf("%p", actual),
		fmt.Sprintf("%p", expected),
		"actual and expected pointers should be equal",
	)
}

func (c *TestContainer) MustNotEqualPointer(expected interface{}, actual interface{}) {
	require.NotEqual(c.t,
		fmt.Sprintf("%p", actual),
		fmt.Sprintf("%p", expected),
		"actual and expected pointers should not be equal",
	)
}
