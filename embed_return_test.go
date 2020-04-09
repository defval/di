package di

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type Config struct {
	Port   int
}

type Server struct {
	Cfg *Config `di:""`
}

func NewConfig() *Config {
	return &Config{Port: 8080}
}

func NewServer() *Server {
	return &Server{}
}

func TestEmbedReturn(t *testing.T) {
	t.Run("embed return success", func(t *testing.T) {
		c, err := New(
			Provide(NewConfig),
			Provide(NewServer),
		)
		if err != nil {
			t.Fatal(err)
		}
		defer c.Cleanup()
		err = c.Invoke(func(server *Server) {
			assert.Equal(t, server.Cfg, &Config{
				Port: 8080,
			})
		})
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("embed return cycle cause error", func(t *testing.T) {
		_, err := New(
			Provide(NewConfig2),
			Provide(NewServer2),
		)
		require.NotNil(t, err)
		f1 := err.Error() == "[*di.Config2 *di.Server2 *di.Config2] cycle detected"
		f2 := err.Error() == "[*di.Server2 *di.Config2 *di.Server2] cycle detected"
		require.True(t, f1 || f2, err.Error())
	})

}

type Config2 struct {
	Port   int
	Server *Server2 `di:""`
}

type Server2 struct {
	Cfg *Config2 `di:""`
}

func NewConfig2() *Config2 {
	return &Config2{Port: 8080}
}

func NewServer2() *Server2 {
	return &Server2{}
}