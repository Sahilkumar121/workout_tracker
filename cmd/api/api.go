package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Sahilkumar121/workout_tracker/internal/config"
	"github.com/Sahilkumar121/workout_tracker/internal/db"
	"github.com/Sahilkumar121/workout_tracker/internal/handler"
	"github.com/Sahilkumar121/workout_tracker/internal/middleware"
)

func main() {

	// slog setup
	loggerJSONHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(loggerJSONHandler)
	slog.SetDefault(logger)

	// load config
	cfg := config.MustLoad()

	// database connection
	pgDB, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("failed database load %v", err)
	}
	fmt.Println("database connect successfully ...")

	// new mux server
	mux := http.NewServeMux()

	// new api handler
	user := handler.NewUserHandler(pgDB, logger, cfg.SecreteKey)
	workout := handler.NewWorkoutHandler(pgDB, logger)

	// api
	mux.HandleFunc("GET /health", handler.CheckHealth)
	mux.HandleFunc("POST /user/register", user.RegisterUser)
	mux.HandleFunc("POST /user/login", user.LoginUser)

	// user api
	mux.Handle("GET /user/me", middleware.Auth(http.HandlerFunc(user.GetProfile)))
	mux.Handle("PATCH /user/update", middleware.Auth(http.HandlerFunc(user.UpdateUser)))

	// workout api
	mux.Handle("POST /user/workout", middleware.Auth(http.HandlerFunc(workout.PostWorkout)))
	mux.Handle("GET /user/workout", middleware.Auth(http.HandlerFunc(workout.GetUserWorkouts)))

	// middleware handler
	requestIDHandler := middleware.RequestID(mux)
	fmt.Println("server is starting ...")
	ser := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      requestIDHandler,
		ReadTimeout:  time.Minute * 5,
		WriteTimeout: time.Minute * 5,
		IdleTimeout:  time.Hour * 1,
	}

	fmt.Printf("server has started running on http://localhost:%s\n", cfg.Port)
	err = ser.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start the server %v", err)
	}
}
