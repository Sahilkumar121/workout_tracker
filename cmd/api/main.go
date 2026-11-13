package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sahilkumar121/workout_tracker/internal/config"
	"github.com/Sahilkumar121/workout_tracker/internal/database"
	"github.com/Sahilkumar121/workout_tracker/internal/handler/exercise"
	"github.com/Sahilkumar121/workout_tracker/internal/repository"
	"github.com/Sahilkumar121/workout_tracker/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

// createServer return http.Server
func createServer(
	cfg config.Config,
	pool *pgxpool.Pool,
	logger *slog.Logger,
) *http.Server {

	mux := http.NewServeMux()

	handleAPI(mux, pool, logger)

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}
}

// runServer return error if any server, or interrupt occur and shutdown gracefully
func runServer(
	ctx context.Context,
	server *http.Server,
	serverShutdownTimeout time.Duration,
) error {

	// server error channel with buffer size of 1
	serverErr := make(chan error, 1)

	// signal error channel with buffer size of 1
	stop := make(chan os.Signal, 1)

	// Notify func is a write channel function
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	go func() {
		log.Println("Server Started on http://localhost:8080/api/v1")
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil || errors.Is(err, http.ErrServerClosed) {
			log.Println("Server failed to start")
			return err
		}

	case sig := <-stop:
		log.Printf("Interrupt signal received %v\n", sig)

	case <-ctx.Done():
		log.Println("Context cancelled...")
	}

	// gracefully shutdown
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		serverShutdownTimeout,
	)
	defer cancel()

	// try gracefully shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("Failed to gracefully shutdown")

		// forcefully close the server
		if closeErr := server.Close(); closeErr != nil {
			log.Println("Forcefully server close failed")
			return errors.Join(err, closeErr)
		}

		return err
	}

	log.Println("Gracefully Server Shutdown...")
	return nil
}

func handleAPI(
	mux *http.ServeMux,
	pool *pgxpool.Pool,
	logger *slog.Logger,
) {

	// repos
	exerciseRepo := repository.NewExerciseRepository(logger, pool)

	// services
	exerciseService := service.NewExerciseServiceHandler(logger, exerciseRepo)

	// handler
	exerciseHandler := exercise.NewExerciseHandler(logger, exerciseService)

	// public api
	mux.HandleFunc("POST /api/v1/exercise", exerciseHandler.PostExercise)
}

func main() {

	// logger
	loggerHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)

	// configuration
	cfg := config.MustLoad()

	// database
	postgresCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	pool, err := database.PostgresqlDB(postgresCtx, cfg.DbUrl)
	if err != nil {
		log.Fatalf("initialization of postgres failed: %v\n", err)
	}

	// create and run server
	server := createServer(cfg, pool, logger)
	if err := runServer(context.Background(), server, 5*time.Second); err != nil {
		log.Fatalf("Error occur %v\n", err)
	}
}
