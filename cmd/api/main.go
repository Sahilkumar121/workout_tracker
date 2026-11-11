package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sahilkumar121/workout_tracker/internal/config"
)

// createServer return http.Server
func createServer(cfg config.Config) *http.Server {

	_ = cfg

	mux := http.NewServeMux()

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

// runServer return error if any server, or interupt occur and shutdown gracefully
func runServer(
	ctx context.Context,
	server *http.Server,
	serverShutdownTimeout time.Duration,
) error {

	// server error channel with buffer size of 1
	serverErr := make(chan error, 1)

	// signal error channel with buffer size of 1
	stop := make(chan os.Signal, 1)

	// Notify func is a write chan function
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	go func() {
		log.Println("Server Started on http://localhost:8080/api/v1")
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil || errors.Is(err, http.ErrServerClosed) {
			log.Println("Server failed to start\n")
			return err
		}

	case sig := <-stop:
		log.Printf("Interupt signal received %v\n", sig)

	case <-ctx.Done():
		log.Println("Context cancelled...")
	}

	// gracefull shutdown
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		serverShutdownTimeout,
	)
	defer cancel()

	// try gracefull shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("Failed to gracefull shutdown\n")

		// forcefully close the server
		if closeErr := server.Close(); closeErr != nil {
			log.Println("Forcefull server close failed\n")
			return errors.Join(err, closeErr)
		}

		return err
	}

	log.Println("Gracefull Server Shutdown...")
	return nil
}

func main() {

	cfg := config.MustLoad()

	// create and run server
	server := createServer(cfg)
	if err := runServer(context.Background(), server, 5*time.Second); err != nil {
		log.Fatalf("Error occur %v\n", err)
	}
}
