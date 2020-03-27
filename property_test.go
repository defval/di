package di

import "testing"

type Config struct {
	Port int
	server *Server
}

type Server struct {
	cfg *Config `di:""`
}

func NewConfig() *Config {
	return &Config{Port:8080}
}

func NewServer() *Server {
	return &Server{}
}

func TestAs(t *testing.T) {
	c, err := New(
		Provide(NewConfig),
		Provide(NewServer),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Invoke(func(server *Server) {
		t.Log(server.cfg)
	})
	if err != nil {
		t.Fatal(err)
	}
	c.Cleanup()
}
