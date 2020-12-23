package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/goava/di"
)

func main() {
	di.SetTracer(&di.StdTracer{})
	// create container
	c, err := di.New(
		di.Provide(NewContext),  // provide application context
		di.Provide(NewServer),   // provide http server
		di.Provide(NewServeMux), // provide http serve mux
		// controllers as []Controller group
		di.Provide(NewOrderController, di.As(new(Controller))),
		di.Provide(NewUserController, di.As(new(Controller))),
	)
	// handle container errors
	if err != nil {
		log.Fatal(err)
	}
	// invoke function
	if err := c.Invoke(StartServer); err != nil {
		log.Fatal(err)
	}
}

// StartServer starts http server.
func StartServer(ctx context.Context, server *http.Server) error {
	log.Println("start server")
	errChan := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()
	select {
	case <-ctx.Done():
		log.Println("stop server")
		return server.Close()
	case err := <-errChan:
		return fmt.Errorf("server error: %s", err)
	}
}

// NewContext creates new application context.
func NewContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
		<-stop
		cancel()
	}()
	return ctx
}

// NewServer creates a http server with provided mux as handler.
func NewServer(mux *http.ServeMux) *http.Server {
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	return server
}

// NewServeMux creates a new http serve mux.
func NewServeMux(controllers []Controller) *http.ServeMux {
	mux := &http.ServeMux{}
	for _, controller := range controllers {
		controller.RegisterRoutes(mux)
	}
	return mux
}

// Controller is an interface that can register its routes.
type Controller interface {
	RegisterRoutes(mux *http.ServeMux)
}

// OrderController is a http controller for orders.
type OrderController struct{}

// NewOrderController creates a auth http controller.
func NewOrderController() *OrderController {
	return &OrderController{}
}

// RegisterRoutes is a Controller interface implementation.
func (a *OrderController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/orders", a.RetrieveOrders)
}

// Retrieve loads orders and writes it to the writer.
func (a *OrderController) RetrieveOrders(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("Orders"))
}

// UserController is a http endpoint for a user.
type UserController struct{}

// NewUserController creates a user http endpoint.
func NewUserController() *UserController {
	return &UserController{}
}

// RegisterRoutes is a Controller interface implementation.
func (e *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", e.RetrieveUsers)
}

// Retrieve loads users and writes it using the writer.
func (e *UserController) RetrieveUsers(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("Users"))
}
