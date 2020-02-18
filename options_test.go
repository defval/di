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
		c := di.New(
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
		require.NoError(t, c.Compile())
		require.Equal(t, loadedServer, server)
		require.Equal(t, loadedServer.Handler, mux)
		require.Equal(t, resolvedServer, server)
	})

	t.Run("provide failed on compile", func(t *testing.T) {
		c := di.New(
			di.Provide(func() {}),
		)
		require.EqualError(t, c.Compile(), "constructor must be a function like func([dep1, dep2, ...]) (<result>, [cleanup, error]), got func()")
	})
}
