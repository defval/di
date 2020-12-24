package di_test

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/goava/di"
)

func init() {
	di.SetTracer(di.StdTracer{})
}

func TestContainer_Provide(t *testing.T) {
	t.Run("simple constructor", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NoError(t, c.Provide(func() *http.Server { return &http.Server{} }))
	})

	t.Run("constructor with cleanup function", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NoError(t, c.Provide(func() (*http.Server, func()) {
			return &http.Server{}, func() {}
		}))
	})

	t.Run("constructor with cleanup and error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NoError(t, c.Provide(func() (*http.Server, func(), error) {
			return &http.Server{}, func() {}, nil
		}))
	})

	t.Run("provide string cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Provide("string")
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": invalid constructor signature, got string")
	})

	t.Run("provide nil cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Provide(nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": invalid constructor signature, got nil")
	})

	t.Run("provide struct pointer cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Provide(&http.Server{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": invalid constructor signature, got *http.Server")
	})

	t.Run("provide constructor without result cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Provide(func() {})
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": invalid constructor signature, got func()")
	})

	t.Run("provide constructor with many resultant types cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		ctor := func() (*http.Server, *http.ServeMux, error) {
			return nil, nil, nil
		}
		err = c.Provide(ctor)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": invalid constructor signature, got func() (*http.Server, *http.ServeMux, error)")
	})

	t.Run("provide constructor with incorrect result error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		ctor := func() (*http.Server, *http.ServeMux) {
			return nil, nil
		}
		err = c.Provide(ctor)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), "invalid constructor signature, got func() (*http.Server, *http.ServeMux)")
	})

	t.Run("provide duplicate not cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		ctor := func() *http.Server { return nil }
		require.NoError(t, c.Provide(ctor))
		require.NoError(t, c.Provide(ctor))
	})

	t.Run("provide as not implemented interface cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		// http server not implement io.Reader interface
		err = c.Provide(func() *http.Server { return nil }, di.As(new(io.Reader)))
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": *http.Server not implement io.Reader")
	})

	t.Run("provide type as several interfaces", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		file := &os.File{}
		require.NoError(t, c.Provide(func() *os.File { return file }, di.As(new(io.Closer), new(io.ReadCloser))))
		var closer io.Closer
		var readCloser io.ReadCloser
		require.NoError(t, c.Resolve(&closer))
		require.NoError(t, c.Resolve(&readCloser))
		require.Equal(t, fmt.Sprintf("%p", closer), fmt.Sprintf("%p", file))
		require.Equal(t, fmt.Sprintf("%p", readCloser), fmt.Sprintf("%p", file))
	})

	t.Run("using not interface type in di.As() cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Provide(func() *http.Server { return nil }, di.As(&http.Server{}))
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": *http.Server: not a pointer to interface")
	})

	t.Run("using nil type in di.As() cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Provide(func() *http.Server { return &http.Server{} }, di.As(nil))
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": nil: not a pointer to interface")
	})
}

func TestContainer_Resolve(t *testing.T) {
	t.Run("resolve into nil cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		err = c.Resolve(nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": target must be a pointer, got nil")
	})

	t.Run("resolve into struct cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		err = c.Resolve(struct{}{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": target must be a pointer, got struct {}")
	})

	t.Run("resolve into string cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		err = c.Resolve("string")
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": target must be a pointer, got string")
	})

	t.Run("resolve with failed build", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Provide(func() (*http.Server, error) {
			return &http.Server{}, fmt.Errorf("server build failed")
		})
		require.NoError(t, err)
		var server *http.Server
		err = c.Resolve(&server)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": *http.Server: server build failed")
	})

	t.Run("resolve returns type that was created in constructor", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		server := &http.Server{}
		require.NoError(t, c.Provide(func() *http.Server { return server }))
		var extracted *http.Server
		require.NoError(t, c.Resolve(&extracted))
		require.Equal(t, fmt.Sprintf("%p", server), fmt.Sprintf("%p", extracted))
	})

	t.Run("resolve same pointer on each resolve", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NoError(t, c.Provide(func() *http.Server { return &http.Server{} }))
		var server1 *http.Server
		require.NoError(t, c.Resolve(&server1))
		var server2 *http.Server
		require.NoError(t, c.Resolve(&server2))
		require.Equal(t, fmt.Sprintf("%p", server1), fmt.Sprintf("%p", server2))
	})

	t.Run("resolve not existing type cause error", func(t *testing.T) {
		c, err := di.New()
		require.NotNil(t, c)
		require.NoError(t, err)
		err = c.Resolve(&http.Server{})
		require.Error(t, err)
		require.True(t, errors.Is(err, di.ErrTypeNotExists))
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": type http.Server not exists in the container")
	})

	t.Run("resolve functions", func(t *testing.T) {
		var result []string
		fn1 := func() { result = append(result, "fn1") }
		fn2 := func() { result = append(result, "fn2") }
		fn3 := func() { result = append(result, "fn3") }
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		type MyFunc func()
		require.NoError(t, c.Provide(func() MyFunc { return fn1 }))
		require.NoError(t, c.Provide(func() MyFunc { return fn2 }))
		require.NoError(t, c.Provide(func() MyFunc { return fn3 }))
		var funcs []MyFunc
		require.NoError(t, c.Resolve(&funcs))
		require.Len(t, funcs, 3)
		for _, fn := range funcs {
			fn()
		}
		require.Equal(t, []string{"fn1", "fn2", "fn3"}, result)
	})

	t.Run("container provided by default", func(t *testing.T) {
		var container *di.Container
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NoError(t, c.Resolve(&container))
		require.Equal(t, fmt.Sprintf("%p", c), fmt.Sprintf("%p", container))
	})

	t.Run("cycle cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		// bool -> int32 -> int64 -> bool
		err = c.Provide(func(int32) bool { return true })
		require.NoError(t, err)
		err = c.Provide(func(int64) int32 { return 0 })
		require.NoError(t, err)
		err = c.Provide(func(bool) int64 { return 0 })
		require.NoError(t, err)
		var b bool
		err = c.Resolve(&b)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": cycle detected") // todo: improve message
	})

	t.Run("resolve not existing dependency type cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NoError(t, c.Provide(func(int) int32 { return 0 }))
		var i int32
		err = c.Resolve(&i)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": int32: type int not exists in the container")
	})

	t.Run("resolve correct argument", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		mux := &http.ServeMux{}
		require.NoError(t, c.Provide(func() *http.ServeMux { return mux }))
		require.NoError(t, c.Provide(func(mux *http.ServeMux) *http.Server {
			return &http.Server{Handler: mux}
		}))
		var server *http.Server
		require.NoError(t, c.Resolve(&server))
		require.Equal(t, fmt.Sprintf("%p", mux), fmt.Sprintf("%p", server.Handler))
	})
}

func TestContainer_Interfaces(t *testing.T) {
	t.Run("resolve interface with several implementations cause error", func(t *testing.T) {
		c, err := di.New(
			di.Provide(func() *http.Server { return &http.Server{} }, di.As(new(io.Closer))),
			di.Provide(func() *os.File { return &os.File{} }, di.As(new(io.Closer))),
		)
		require.NoError(t, err)
		var closer io.Closer
		err = c.Resolve(&closer)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": multiple definitions of io.Closer, maybe you need to use group type: []io.Closer")
	})

	t.Run("resolve constructor interface argument", func(t *testing.T) {
		mux := &http.ServeMux{}
		c, err := di.New(
			di.Provide(func() *http.ServeMux { return mux }, di.As(new(http.Handler))),
			di.Provide(func(handler http.Handler) *http.Server { return &http.Server{Handler: handler} }),
		)
		require.NoError(t, err)
		var handler http.Handler
		err = c.Resolve(&handler)
		require.NoError(t, err)
		var server *http.Server
		err = c.Resolve(&server)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%p", mux), fmt.Sprintf("%p", server.Handler))
	})

	t.Run("resolve not existing unnamed definition with named", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NoError(t, c.Provide(http.NewServeMux, di.WithName("two")))
		require.NoError(t, c.Provide(http.NewServeMux, di.WithName("three")))
		var mux *http.ServeMux
		err = c.Resolve(&mux)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": multiple definitions of *http.ServeMux, maybe you need to use group type: []*http.ServeMux")
	})

	t.Run("resolve same pointer on resolve", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NoError(t, c.Provide(func() *http.ServeMux { return &http.ServeMux{} }, di.As(new(http.Handler))))
		var server *http.ServeMux
		require.NoError(t, c.Resolve(&server))
		var handler http.Handler
		require.NoError(t, c.Resolve(&handler))
		require.Equal(t, fmt.Sprintf("%p", server), fmt.Sprintf("%p", handler))
	})
}

func TestContainer_Groups(t *testing.T) {
	t.Run("resolve multiple type instances as slice of type", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		conn1 := &net.TCPConn{}
		conn2 := &net.TCPConn{}
		require.NoError(t, c.Provide(func() *net.TCPConn { return conn1 }))
		require.NoError(t, c.Provide(func() *net.TCPConn { return conn2 }))
		var conns []*net.TCPConn
		require.NoError(t, c.Resolve(&conns))
		require.Len(t, conns, 2)
		require.Equal(t, fmt.Sprintf("%p", conn1), fmt.Sprintf("%p", conns[0]))
		require.Equal(t, fmt.Sprintf("%p", conn2), fmt.Sprintf("%p", conns[1]))
	})

	t.Run("resolve not specific type of group cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		conn1 := &net.TCPConn{}
		conn2 := &net.TCPConn{}
		require.NoError(t, c.Provide(func() *net.TCPConn { return conn1 }))
		require.NoError(t, c.Provide(func() *net.TCPConn { return conn2 }))
		var conn *net.TCPConn
		err = c.Resolve(&conn)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": multiple definitions of *net.TCPConn, maybe you need to use group type: []*net.TCPConn")
	})

	t.Run("resolve group of interface", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		server := &http.Server{}
		file := &os.File{}
		require.NoError(t, c.Provide(func() *http.Server { return server }, di.As(new(io.Closer))))
		require.NoError(t, c.Provide(func() *os.File { return file }, di.As(new(io.Closer))))
		var closers []io.Closer
		require.NoError(t, c.Resolve(&closers))
		require.Len(t, closers, 2)
		require.Equal(t, fmt.Sprintf("%p", server), fmt.Sprintf("%p", closers[0]))
		require.Equal(t, fmt.Sprintf("%p", file), fmt.Sprintf("%p", closers[1]))
	})

	t.Run("group updates on provide", func(t *testing.T) {
		var result []string
		fn1 := func() { result = append(result, "fn1") }
		fn2 := func() { result = append(result, "fn2") }
		fn3 := func() { result = append(result, "fn3") }
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		type MyFunc func()
		var funcs []MyFunc
		require.NoError(t, c.Provide(func() MyFunc { return fn1 }))
		require.NoError(t, c.Resolve(&funcs))
		require.Len(t, funcs, 1)
		require.NoError(t, c.Provide(func() MyFunc { return fn2 }))
		require.NoError(t, c.Resolve(&funcs))
		require.Len(t, funcs, 2)
		require.NoError(t, c.Provide(func() MyFunc { return fn3 }))
		require.NoError(t, c.Resolve(&funcs))
		require.Len(t, funcs, 3)
	})

	t.Run("resolve one interface from group of type", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		conn1 := &net.TCPConn{}
		conn2 := &net.TCPConn{}
		require.NoError(t, c.Provide(func() *net.TCPConn { return conn1 }))
		require.NoError(t, c.Provide(func() *net.TCPConn { return conn2 }, di.As(new(net.Conn))))
		var conn net.Conn
		require.NoError(t, c.Resolve(&conn))
		require.Equal(t, fmt.Sprintf("%p", conn), fmt.Sprintf("%p", conn))
	})
}

func TestContainer_Iterate(t *testing.T) {
	t.Run("iterate over nil causes error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		err = c.Iterate(nil, func(tags di.Tags, loader di.ValueFunc) error {
			return nil
		})
		require.EqualError(t, err, "target must be a pointer, got nil")
	})
	t.Run("iterate over struct causes error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		err = c.Iterate(http.ServeMux{}, func(tags di.Tags, loader di.ValueFunc) error {
			return nil
		})
		require.EqualError(t, err, "target must be a pointer, got http.ServeMux")
	})
	t.Run("iterate over struct causes error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NoError(t, c.Provide(func() http.ServeMux { return http.ServeMux{} }))
		err = c.Iterate(&http.ServeMux{}, func(tags di.Tags, loader di.ValueFunc) error {
			return nil
		})
		require.EqualError(t, err, "iteration can be used with groups only")
	})
	t.Run("iterates over instances", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		conn1 := &net.TCPConn{}
		conn2 := &net.TCPConn{}
		require.NoError(t, c.Provide(func() *net.TCPConn { return conn1 }))
		require.NoError(t, c.Provide(func() *net.TCPConn { return conn2 }))
		var iterates []*net.TCPConn
		var conn []*net.TCPConn
		iterFn := func(tags di.Tags, loader di.ValueFunc) error {
			i, err := loader()
			if err != nil {
				return err
			}
			iterates = append(iterates, i.(*net.TCPConn))
			return nil
		}
		err = c.Iterate(&conn, iterFn)
		require.NoError(t, err)
		require.Len(t, iterates, 2)
		require.Equal(t, iterates[0], conn1)
		require.Equal(t, iterates[1], conn2)
	})

	t.Run("iterates over tagged instances", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		conn1 := &net.TCPConn{}
		conn2 := &net.TCPConn{}
		require.NoError(t, c.Provide(func() *net.TCPConn { return conn1 }, di.Tags{"conn": "tcp1"}))
		require.NoError(t, c.Provide(func() *net.TCPConn { return conn2 }, di.Tags{"conn": "tcp2"}))
		require.NoError(t, c.Provide(func() *net.TCPConn { return &net.TCPConn{} }))
		var iterates []*net.TCPConn
		var all []di.Tags
		var conn []*net.TCPConn
		iterFn := func(tags di.Tags, loader di.ValueFunc) error {
			all = append(all, tags)
			i, err := loader()
			if err != nil {
				return err
			}
			iterates = append(iterates, i.(*net.TCPConn))
			return nil
		}
		err = c.Iterate(&conn, iterFn, di.Tags{"conn": "*"})
		require.NoError(t, err)
		require.Len(t, iterates, 2)
		require.Equal(t, conn1, iterates[0])
		require.Equal(t, conn2, iterates[1])
		require.Equal(t, []di.Tags{
			{"conn": "tcp1"},
			{"conn": "tcp2"},
		}, all)
	})

	t.Run("iterates over instances with errors", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		conn1 := &net.TCPConn{}
		conn2 := &net.TCPConn{}
		require.NoError(t, c.Provide(func() *net.TCPConn { return conn1 }))
		require.NoError(t, c.Provide(func() (*net.TCPConn, error) { return conn2, fmt.Errorf("tcp conn 2 error") }))
		var iterates []*net.TCPConn
		var conn []*net.TCPConn
		iterFn := func(tags di.Tags, loader di.ValueFunc) error {
			i, err := loader()
			if err != nil {
				return err
			}
			iterates = append(iterates, i.(*net.TCPConn))
			return nil
		}
		err = c.Iterate(&conn, iterFn)
		require.EqualError(t, err, "[]*net.TCPConn with index 1 failed: tcp conn 2 error")
	})
}

func TestContainer_Tags(t *testing.T) {
	t.Run("resolve named definition", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		first := &http.Server{}
		second := &http.Server{}
		err = c.Provide(func() *http.Server { return first }, di.WithName("first"))
		require.NoError(t, err)
		err = c.Provide(func() *http.Server { return second }, di.WithName("second"))
		var extracted *http.Server
		err = c.Resolve(&extracted)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": multiple definitions of *http.Server, maybe you need to use group type: []*http.Server")
		err = c.Resolve(&extracted, di.Name("first"))
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%p", first), fmt.Sprintf("%p", extracted))
		err = c.Resolve(&extracted, di.Name("second"))
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%p", second), fmt.Sprintf("%p", extracted))
	})

	t.Run("resolve single instance of group without specifying tags cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NoError(t, c.Provide(http.NewServeMux, di.WithName("first")))
		var mux *http.ServeMux
		require.NoError(t, c.Resolve(&mux))
		require.NoError(t, c.Provide(http.NewServeMux, di.WithName("second")))
		err = c.Resolve(&mux)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": multiple definitions of *http.ServeMux, maybe you need to use group type: []*http.ServeMux")
	})

	t.Run("resolve not found by tags instance cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NoError(t, c.Provide(http.NewServeMux, di.WithName("first")))
		require.NoError(t, c.Provide(http.NewServeMux, di.WithName("second")))
		var mux *http.ServeMux
		err = c.Resolve(&mux, di.Name("unknown"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": type *http.ServeMux[name:unknown] not exists")
	})

	t.Run("provide duplication of named definition", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NoError(t, c.Provide(http.NewServeMux, di.WithName("first")))
		err = c.Provide(http.NewServeMux, di.WithName("first"))
		require.NoError(t, err)
	})

	t.Run("resolve existing unnamed definition with named", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NoError(t, c.Provide(http.NewServeMux))
		require.NoError(t, c.Provide(http.NewServeMux, di.WithName("two")))
		require.NoError(t, c.Provide(http.NewServeMux, di.WithName("three")))
		var mux *http.ServeMux
		err = c.Resolve(&mux)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), "multiple definitions of *http.ServeMux, maybe you need to use group type: []*http.ServeMux")
		require.NoError(t, c.Resolve(&mux, di.Name("two")))
		require.NoError(t, c.Resolve(&mux, di.Name("three")))
	})

	t.Run("resolve instances with same tag", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NoError(t, c.Provide(http.NewServeMux))
		require.NoError(t, c.Provide(http.NewServeMux, di.Tags{"tag": "the_same"}))
		require.NoError(t, c.Provide(http.NewServeMux, di.Tags{"tag": "the_same"}))
		var muxs []*http.ServeMux
		err = c.Resolve(&muxs, di.Tags{"tag": "the_same"})
		require.NoError(t, err)
		require.Len(t, muxs, 2)
	})

	t.Run("resolve all instances with tag", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NoError(t, c.Provide(http.NewServeMux))
		require.NoError(t, c.Provide(http.NewServeMux, di.Tags{"server": "one"}))
		require.NoError(t, c.Provide(http.NewServeMux, di.Tags{"server": "two"}))
		var muxs []*http.ServeMux
		err = c.Resolve(&muxs, di.Tags{"server": "*"})
		require.NoError(t, err)
		require.Len(t, muxs, 2)
	})

	t.Run("resolve all instances with several tags", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NoError(t, c.Provide(http.NewServeMux))
		require.NoError(t, c.Provide(http.NewServeMux, di.Tags{"server": "one"}))
		require.NoError(t, c.Provide(http.NewServeMux, di.Tags{"server": "one", "http": "one"}))
		require.NoError(t, c.Provide(http.NewServeMux, di.Tags{"server": "two", "http": "two"}))
		var muxs []*http.ServeMux
		err = c.Resolve(&muxs, di.Tags{"server": "*", "http": "*"})
		require.NoError(t, err)
		require.Len(t, muxs, 2)
	})

	t.Run("provide type with tags", func(t *testing.T) {
		type Server struct {
			di.Tags `http:"true" server:"true"`
		}
		var s *Server
		_, err := di.New(
			di.Provide(func() *Server { return &Server{} }),
			di.Resolve(&s, di.Tags{"http": "true", "server": "true"}),
		)
		require.NoError(t, err)
	})
}

func TestContainer_Group(t *testing.T) {
	t.Run("resolve group argument", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		server := &http.Server{}
		file := &os.File{}
		require.NoError(t, c.Provide(func() *http.Server { return server }, di.As(new(io.Closer))))
		require.NoError(t, c.Provide(func() *os.File { return file }, di.As(new(io.Closer))))
		type Closers []io.Closer
		require.NoError(t, c.Provide(func(closers []io.Closer) Closers { return closers }))
		var closers Closers
		require.NoError(t, c.Resolve(&closers))
		require.Equal(t, fmt.Sprintf("%p", server), fmt.Sprintf("%p", closers[0]))
		require.Equal(t, fmt.Sprintf("%p", file), fmt.Sprintf("%p", closers[1]))
	})

	t.Run("incorrect signature", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Invoke(func() *http.Server { return &http.Server{} })
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": invalid invocation signature, got func() *http.Server")
	})
}

func TestContainer_Invoke(t *testing.T) {
	t.Run("invoke nil", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Invoke(nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": invalid invocation signature, got nil")
	})
	t.Run("invoke invalid function", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Invoke(func() *http.Server { return &http.Server{} })
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": invalid invocation signature, got func() *http.Server")
	})
	t.Run("invocation function with not provided dependency cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Invoke(func(server *http.Server) {})
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": type *http.Server not exists in the container")
	})

	t.Run("invoke with nil error must be called", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		var invokeCalled bool
		err = c.Invoke(func() error {
			invokeCalled = true
			return nil
		})
		require.NoError(t, err)
		require.True(t, invokeCalled)
	})

	t.Run("resolve dependencies in invoke", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		server := &http.Server{}
		called := false
		require.NoError(t, c.Provide(func() *http.Server { return server }))
		err = c.Invoke(func(in *http.Server) {
			called = true
			require.Equal(t, fmt.Sprintf("%p", server), fmt.Sprintf("%p", in))
		})
		require.NoError(t, err)
		require.True(t, called)
	})

	t.Run("invoke return error as is", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Invoke(func() error { return fmt.Errorf("invoke error") })
		require.EqualError(t, err, "invoke error")
	})

	t.Run("cycle cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		// bool -> int32 -> int64 -> bool
		err = c.Provide(func(int32) bool { return true })
		require.NoError(t, err)
		err = c.Provide(func(int64) int32 { return 0 })
		require.NoError(t, err)
		err = c.Provide(func(bool) int64 { return 0 })
		require.NoError(t, err)
		err = c.Invoke(func(bool) {})
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": cycle detected") // todo: improve message
	})
}

func TestContainer_Has(t *testing.T) {
	t.Run("exists nil returns false", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		has, err := c.Has(nil)
		require.EqualError(t, err, "target must be a pointer, got nil")
		require.False(t, has)
	})

	t.Run("exists return true if type exists", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NoError(t, c.Provide(func() *http.Server { return &http.Server{} }))
		var server *http.Server
		has, err := c.Has(&server)
		require.NoError(t, err)
		require.True(t, has)
	})

	t.Run("exists return false if type not exists", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		var server *http.Server
		has, err := c.Has(&server)
		require.NoError(t, err)
		require.False(t, has)
	})

	t.Run("exists interface", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NoError(t, c.Provide(func() *http.Server { return &http.Server{} }, di.As(new(io.Closer))))
		var server io.Closer
		has, err := c.Has(&server)
		require.NoError(t, err)
		require.True(t, has)
	})

	t.Run("exists named provider", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		err = c.Provide(func() *http.Server { return &http.Server{} }, di.Tags{"name": "server"})
		require.NoError(t, err)
		var server *http.Server
		has, err := c.Has(&server, di.Tags{"name": "server"})
		require.NoError(t, err)
		require.True(t, has)
	})

	t.Run("type exists but no possible to build returns true", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NoError(t, c.Provide(func(b bool) *http.Server { return &http.Server{} }))
		var server *http.Server
		has, err := c.Has(&server)
		require.EqualError(t, err, "*http.Server: type bool not exists in the container")
		require.False(t, has)
	})
}

func TestContainer_Inject(t *testing.T) {
	t.Run("constructor with injectable pointer", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		type InjectableType struct {
			di.Inject
			Mux *http.ServeMux
		}
		mux := &http.ServeMux{}
		require.NoError(t, c.Provide(func() *http.ServeMux { return mux }))
		require.NoError(t, c.Provide(func() *InjectableType { return &InjectableType{} }))
		var result *InjectableType
		require.NoError(t, c.Resolve(&result))
		require.NotNil(t, result.Mux)
		require.Equal(t, fmt.Sprintf("%p", mux), fmt.Sprintf("%p", result.Mux))
	})

	// todo: https://github.com/goava/di/issues/29
	//t.Run("constructor with injectable embed pointer", func(t *testing.T) {
	//	c, err := di.New()
	//	require.NoError(t, err)
	//	type InjectableType struct {
	//		di.Inject
	//		*http.ServeMux
	//	}
	//	mux := &http.ServeMux{}
	//	require.NoError(t, c.Provide(func() *http.ServeMux { return mux }))
	//	require.NoError(t, c.Provide(func() *InjectableType { return &InjectableType{} }))
	//	var result *InjectableType
	//	require.NoError(t, c.Resolve(&result))
	//	require.NotNil(t, result.ServeMux)
	//	require.Equal(t, fmt.Sprintf("%p", mux), fmt.Sprintf("%p", result.ServeMux))
	//})

	t.Run("provide injectable struct cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		type InjectableType struct {
			di.Inject
			Mux *http.ServeMux
		}
		mux := &http.ServeMux{}
		require.NoError(t, c.Provide(func() *http.ServeMux { return mux }))
		err = c.Provide(func() InjectableType { return InjectableType{} })
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": di.Inject not supported for unaddressable result of constructor, use *di_test.InjectableType instead")
	})

	t.Run("container resolve injectable parameter", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		type Parameters struct {
			di.Inject
			Server *http.Server
			File   *os.File
		}
		server := &http.Server{}
		file := &os.File{}
		require.NoError(t, c.Provide(func() *http.Server { return server }))
		require.NoError(t, c.Provide(func() *os.File { return file }))
		type Result struct {
			server *http.Server
			file   *os.File
		}
		require.NoError(t, c.Provide(func(params Parameters) *Result {
			return &Result{params.Server, params.File}
		}))
		var extracted *Result
		require.NoError(t, c.Resolve(&extracted))
		require.Equal(t, fmt.Sprintf("%p", server), fmt.Sprintf("%p", extracted.server))
		require.Equal(t, fmt.Sprintf("%p", file), fmt.Sprintf("%p", extracted.file))
	})

	t.Run("not existing injectable field cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		type InjectableType struct {
			di.Inject
			Mux *http.ServeMux
		}
		require.NoError(t, c.Provide(func() *InjectableType { return &InjectableType{} }))
		var result *InjectableType
		err = c.Resolve(&result)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": *di_test.InjectableType: type *http.ServeMux not exists in the container")
	})

	t.Run("not existing and optional field set to nil", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		type InjectableType struct {
			di.Inject
			Mux *http.ServeMux `optional:"true"`
		}
		require.NoError(t, c.Provide(func() *InjectableType { return &InjectableType{} }))
		var result *InjectableType
		require.NoError(t, c.Resolve(&result))
		require.Nil(t, result.Mux)
	})

	t.Run("nested injectable field resolved correctly", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		type NestedInjectableType struct {
			di.Inject
			Mux *http.ServeMux
		}
		type InjectableType struct {
			di.Inject
			Nested NestedInjectableType
		}
		mux := &http.ServeMux{}
		require.NoError(t, c.Provide(func() *InjectableType { return &InjectableType{} }))
		require.NoError(t, c.Provide(func() *http.ServeMux { return mux }))
		var result *InjectableType
		require.NoError(t, c.Resolve(&result))
		require.NotNil(t, result.Nested.Mux)
		require.Equal(t, fmt.Sprintf("%p", mux), fmt.Sprintf("%p", result.Nested.Mux))
		var nit NestedInjectableType
		require.NoError(t, c.Resolve(&nit))
		require.NotNil(t, nit)
		require.Equal(t, fmt.Sprintf("%p", mux), fmt.Sprintf("%p", nit.Mux))
	})

	t.Run("optional parameter may be nil", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		type Parameter struct {
			di.Inject
			Server *http.Server `optional:"true"`
		}
		type Result struct {
			server *http.Server
		}
		require.NoError(t, c.Provide(func(params Parameter) *Result { return &Result{server: params.Server} }))
		var extracted *Result
		require.NoError(t, c.Resolve(&extracted))
		require.Nil(t, extracted.server)
	})

	t.Run("resolve group in params", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)

		type Fn func()
		type Params struct {
			di.Inject
			Handlers []Fn `optional:"true"`
		}
		require.NoError(t, c.Provide(func() Fn { return func() {} }))
		require.NoError(t, c.Provide(func() Fn { return func() {} }))
		require.NoError(t, c.Provide(func(params Params) bool {
			return len(params.Handlers) == 2
		}))
		var extracted bool
		require.NoError(t, c.Resolve(&extracted))
		require.True(t, extracted)
	})

	t.Run("optional group may be nil", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		type Params struct {
			di.Inject
			Handlers []http.Handler `optional:"true"`
		}
		require.NoError(t, c.Provide(func(params Params) bool {
			return params.Handlers == nil
		}))
		var extracted bool
		require.NoError(t, c.Resolve(&extracted))
		require.True(t, extracted)
	})

	t.Run("skip private and skip tagged fields", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		type InjectableParameter struct {
			di.Inject
			private []http.Handler
			Addrs   []net.Addr     `optional:"true"`
			Skipped *http.ServeMux `skip:"true"`
		}
		type InjectableType struct {
			di.Inject
			private []http.Handler
			Addrs   []net.Addr `optional:"true"`
		}
		require.NoError(t, c.Provide(func(param InjectableParameter) bool {
			return param.Addrs == nil
		}))
		require.NoError(t, c.Provide(func() *InjectableType { return &InjectableType{} }))
		var extracted bool
		require.NoError(t, c.Resolve(&extracted))
		require.True(t, extracted)
		var result *InjectableType
		require.NoError(t, c.Resolve(&result))
	})

	t.Run("resolving not provided injectable cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		type Parameter struct {
			di.Inject
			Server *http.Server
		}
		var p Parameter
		err = c.Resolve(&p)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": di_test.Parameter: type *http.Server not exists in the container")
	})

	t.Run("resolving provided injectable as interface with dependency", func(t *testing.T) {
		type InjectableType struct {
			di.Inject
			Server *http.Server
		}
		ctor := func() *InjectableType {
			return &InjectableType{}
		}
		server := &http.Server{}
		c, err := di.New(
			di.Provide(func() *http.Server { return server }),
			di.Provide(ctor, di.As(new(di.Interface))),
		)
		require.NoError(t, err)
		var b di.Interface
		err = c.Resolve(&b)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%p", server), fmt.Sprintf("%p", b.(*InjectableType).Server))
	})

	t.Run("resolving provided injectable as interface without dependency cause error", func(t *testing.T) {
		type InjectableType struct {
			di.Inject
			Server *http.Server
		}
		ctor := func() *InjectableType {
			return &InjectableType{}
		}
		c, err := di.New(
			di.Provide(ctor, di.As(new(di.Interface))),
		)
		require.NoError(t, err)
		var b di.Interface
		err = c.Resolve(&b)
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": di.Interface: type *http.Server not exists in the container")
	})

	t.Run("invoke with inject dependency struct", func(t *testing.T) {
		type InjectableParam struct {
			di.Inject
			Mux *http.ServeMux
		}
		c, err := di.New()
		require.NoError(t, err)
		mux := http.NewServeMux()
		require.NoError(t, c.Provide(func() *http.ServeMux { return mux }))
		err = c.Invoke(func(params InjectableParam) {
			require.Equal(t, fmt.Sprintf("%p", mux), fmt.Sprintf("%p", params.Mux))
		})
		require.NoError(t, err)
	})

	t.Run("invoke with inject dependency pointer", func(t *testing.T) {
		type InjectableParam struct {
			di.Inject
			Mux *http.ServeMux
		}
		c, err := di.New()
		require.NoError(t, err)
		mux := http.NewServeMux()
		require.NoError(t, c.Provide(func() *http.ServeMux { return mux }))
		err = c.Invoke(func(params *InjectableParam) {})
		require.Error(t, err)
		require.Contains(t, err.Error(), "container_test.go:")
		require.Contains(t, err.Error(), ": inject *di_test.InjectableParam fields not supported, use di_test.InjectableParam")
	})
}

func TestContainer_Cleanup(t *testing.T) {
	t.Run("called", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		var cleanupCalled bool
		require.NoError(t, c.Provide(func() (*http.Server, func()) {
			return &http.Server{}, func() { cleanupCalled = true }
		}))
		var extracted *http.Server
		require.NoError(t, c.Resolve(&extracted))
		c.Cleanup()
		require.True(t, cleanupCalled)
	})

	t.Run("correct order", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		var cleanupCalls []string
		require.NoError(t, c.Provide(func(handler http.Handler) (*http.Server, func()) {
			return &http.Server{Handler: handler}, func() { cleanupCalls = append(cleanupCalls, "server") }
		}))
		require.NoError(t, c.Provide(func() (*http.ServeMux, func()) {
			return &http.ServeMux{}, func() { cleanupCalls = append(cleanupCalls, "mux") }
		}, di.As(new(http.Handler))))
		var server *http.Server
		require.NoError(t, c.Resolve(&server))
		c.Cleanup()
		require.Equal(t, []string{"server", "mux"}, cleanupCalls)
	})
}
