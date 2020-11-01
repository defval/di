package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	orders := NewOrderController()
	users := NewUserController()
	mux := NewServeMux()
	mux.HandleFunc("/orders", orders.RetrieveOrders)
	mux.HandleFunc("/users", users.RetrieveUsers)
	server := NewServer(mux)
	log.Println("start server")
	errChan := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
		<-stop
		cancel()
	}()
	select {
	case <-ctx.Done():
		log.Println("stop server")
		if err := server.Close(); err != nil {
			log.Fatal(err)
		}
	case err := <-errChan:
		log.Fatal(err)
	}
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
func NewServeMux() *http.ServeMux {
	return &http.ServeMux{}
}

// OrderController is a http controller for orders.
type OrderController struct{}

// NewOrderController creates a auth http controller.
func NewOrderController() *OrderController {
	return &OrderController{}
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

// Retrieve loads users and writes it using the writer.
func (e *UserController) RetrieveUsers(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("Users"))
}
