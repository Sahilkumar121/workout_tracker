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
	"github.com/Sahilkumar121/workout_tracker/internal/handler/user"
	"github.com/Sahilkumar121/workout_tracker/internal/middleware"
	"github.com/Sahilkumar121/workout_tracker/internal/repository"
	"github.com/Sahilkumar121/workout_tracker/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/time/rate"
)

// createServer return http.Server
func createServer(
	cfg config.Config,
	pool *pgxpool.Pool,
	logger *slog.Logger,
) *http.Server {

	mux := http.NewServeMux()
	rateLimiter := middleware.NewIPRateLimiter(rate.Limit(5), 10, 100_000, 10*time.Minute)

	handleAPI(mux, pool, logger, cfg)

	requestID := middleware.RequestId(mux)
	rateLimitMiddleware := rateLimiter.RateLimitMiddleware(requestID)
	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: rateLimitMiddleware,
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
	cfg config.Config,
) {

	// repos
	exerciseRepo := repository.NewExerciseRepository(pool)
	userRepo := repository.NewUserRepository(pool)

	// services
	exerciseService := service.NewExerciseServiceHandler(exerciseRepo)
	userService := service.NewUserService(userRepo)

	// handler
	exerciseHandler := exercise.NewExerciseHandler(logger, exerciseService)
	userHandler := user.NewUserHandler(logger, userService, []byte(cfg.SecretKey))

	// public api

	// user
	mux.HandleFunc("POST /api/v1/user/register", userHandler.RegisterUser)
	mux.HandleFunc("POST /api/v1/user/login", userHandler.LoginUser)
	mux.Handle(
		"GET /api/v1/user/me",
		middleware.AuthenticateUser(
			[]byte(cfg.SecretKey),
			http.HandlerFunc(userHandler.Me),
		),
	)

	// exercise
	mux.HandleFunc("POST /api/v1/exercise", exerciseHandler.PostExercise)
	mux.HandleFunc("GET /api/v1/exercise", exerciseHandler.GetExercises)
	mux.HandleFunc("GET /api/v1/exercise/{id}", exerciseHandler.GetExerciseById)
	mux.HandleFunc("PATCH /api/v1/exercise/{id}", exerciseHandler.PatchExerciseById)
	mux.HandleFunc("DELETE /api/v1/exercise/{id}", exerciseHandler.DeleteExerciseById)
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
