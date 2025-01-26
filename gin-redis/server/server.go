// Package server manages the HTTP server setup and graceful shutdown.
// It initializes the server, starts it, and ensures it shuts down properly when receiving interrupt signals.
package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/keploy/gin-redis/routes"
)

func Init() {
	time.Sleep(2 * time.Second)
	r := routes.NewRouter()
	port := "3001"
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()
	GracefulShutdown(srv)
}
func GracefulShutdown(srv *http.Server) {
	stopper := make(chan os.Signal, 1)
	// listens for interrupt and SIGTERM signal
	signal.Notify(stopper, syscall.SIGINT, syscall.SIGTERM)
	<-stopper
	fmt.Println("Shutting down gracefully...")
	// Create a deadline for the graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server exiting")
}
