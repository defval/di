package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/defval/di"
)

// Configuration
type Configuration struct {
	ConnectionType string
}

// NewConfiguration creates new configuration.
func NewConfiguration() *Configuration {
	c := &Configuration{
		ConnectionType: "tcp",
	}
	if typ, ok := os.LookupEnv("CONNECTION_TYPE"); ok {
		c.ConnectionType = typ
	}
	return c
}

// NewTCPConnection creates tcp connection
func NewTCPConn() *net.TCPConn {
	return &net.TCPConn{}
}

// NewUDPConn creates udp connection
func NewUDPConn() *net.UDPConn {
	return &net.UDPConn{}
}

// ProvideConfiguredConnection
func ProvideConfiguredConnection(conf *Configuration, container *di.Container) error {
	switch conf.ConnectionType {
	case "tcp":
		return container.Provide(NewTCPConn, di.As(new(net.Conn)))
	case "udp":
		return container.Provide(NewUDPConn, di.As(new(net.Conn)))
	}
	return errors.New("unknown connection type")
}

func main() {
	c, err := di.New(
		di.Provide(NewConfiguration),
		di.Invoke(ProvideConfiguredConnection),
	)
	if err != nil {
		log.Fatalln(err)
	}
	var conn net.Conn
	if err := c.Resolve(&conn); err != nil {
		log.Fatalln(err)
	}
	switch conn.(type) {
	case *net.TCPConn:
		fmt.Println("Provided connection: TCP")
	case *net.UDPConn:
		fmt.Println("Provided connection: UDP")
	}
}
