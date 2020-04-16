package di_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/goava/di"
)

func TestOptions(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		var loadedServer *http.Server
		var resolvedServer *http.Server
		server := &http.Server{}
		mux := &http.ServeMux{}
		c, err := di.New(
			di.Options(
				di.Provide(func(handler http.Handler) *http.Server {
					server.Handler = handler
					return server
				}),
				di.Provide(func() *http.ServeMux {
					return mux
				}, di.As(new(http.Handler))),
				di.Invoke(func(server *http.Server) {
					loadedServer = server
				}),
				di.Resolve(&resolvedServer),
			),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.Equal(t, loadedServer, server)
		require.Equal(t, loadedServer.Handler, mux)
		require.Equal(t, resolvedServer, server)
	})

	t.Run("provide failed", func(t *testing.T) {
		c, err := di.New(
			di.Provide(func() {}),
		)
		require.Nil(t, c)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "di.Provide(..) failed:")
		require.Contains(t, err.Error(), "options_test.go:42: constructor must be a function like func([dep1, dep2, ...]) (<result>, [cleanup, error]), got func()")
	})

	t.Run("invoke failed", func(t *testing.T) {
		c, err := di.New(
			di.Invoke(func(string2 string) {}),
		)
		require.Nil(t, c)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "di.Invoke(..) failed:")
		require.Contains(t, err.Error(), "options_test.go:52: resolve invocation (github.com/goava/di_test.TestOptions.func3.1): string: not exists in container")
	})
}
