package di_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defval/di"
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
		require.Contains(t, err.Error(), "options_test.go:")
		require.Contains(t, err.Error(), ": invalid constructor signature, got func()")
	})

	t.Run("invoke failed", func(t *testing.T) {
		c, err := di.New(
			di.Invoke(func(string2 string) {}),
		)
		require.Nil(t, c)
		require.Error(t, err)
		require.Contains(t, err.Error(), "options_test.go:")
		require.Contains(t, err.Error(), ": type string not exists in the container")
	})

	t.Run("invoke error return as is if not internal error", func(t *testing.T) {
		var myError = errors.New("my error")
		_, err := di.New(
			di.Invoke(func() error {
				return myError
			}),
		)
		require.True(t, err == myError)
	})

	t.Run("resolve failed", func(t *testing.T) {
		_, err := di.New(
			di.Resolve(func() {}),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "options_test.go:")
		require.Contains(t, err.Error(), ": target must be a pointer, got func()")
	})
}
