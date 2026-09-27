package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Sahilkumar121/workout_tracker/internal/helper"
	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"github.com/Sahilkumar121/workout_tracker/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkoutHandler struct {
	Logger   *slog.Logger
	DB       *pgxpool.Pool
	Validate *validator.Validate
}

func NewWorkoutHandler(db *pgxpool.Pool, logger *slog.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		Logger:   logger,
		DB:       db,
		Validate: validator.New(),
	}
}

func (h *WorkoutHandler) PostWorkout(w http.ResponseWriter, r *http.Request) {

	requestID := middleware.RequestIDFromContext(r.Context())
	userID := middleware.AuthUserIDFromContext(r.Context())

	var requestData WorkoutRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		h.Logger.Error("failed decode request", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusBadRequest, "incorrect request data", httpx.CodeInvalidInput)
		return
	}

	err = h.Validate.Struct(requestData)
	if err != nil {
		h.Logger.Error("invalid data request", "err", err, "requestId", requestID)
		httpx.Error(w, http.StatusBadRequest, "please provide proper information", httpx.CodeInvalidInput)
		return
	}

	if requestData.WorkoutStartTime == nil {
		startTime := time.Now()
		requestData.WorkoutStartTime = &startTime
	}

	_, err = h.DB.Exec(r.Context(), `insert into workouts (workout_name, workout_start, user_id) values ($1,$2,$3)`, requestData.WorkoutName, requestData.WorkoutStartTime, userID)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			h.Logger.Error("failed create workout (duplicate)", "err", pgxErr, "requestID", requestID)
			httpx.Error(w, http.StatusConflict, "workout already exist", httpx.CodeConflict)
			return
		}
		h.Logger.Error("failed database operation", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong, try again later", httpx.CodeInternalServer)
		return
	}

	helper.SetResponse(w, http.StatusCreated)
}

func (h *WorkoutHandler) GetUserWorkouts(w http.ResponseWriter, r *http.Request) {

	requestID := middleware.RequestIDFromContext(r.Context())
	userID := middleware.AuthUserIDFromContext(r.Context())

	workoutName := r.URL.Query().Get("name")
	dateStr := r.URL.Query().Get("date")

	query := `select id, workout_name, workout_start from workouts where user_id=$1`

	args := []any{userID}
	arrgCount := 1

	if workoutName != "" {
		arrgCount++
		query += fmt.Sprintf(` and workout_name ilike '%%' || $%d || '%%'`, arrgCount)
		args = append(args, workoutName)
	}

	if dateStr != "" {
		arrgCount++
		query += fmt.Sprintf(` and workout_start::date=$%d`, arrgCount)
		args = append(args, dateStr)
	}

	query += ` order by workout_start desc`

	rows, err := h.DB.Query(r.Context(), query, args...)
	if err != nil {
		h.Logger.Error("failed query workout", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong, try again later", httpx.CodeInternalServer)
		return
	}

	collection, err := pgx.CollectRows(rows, pgx.RowToStructByName[WorkoutResponse])
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalServer)
		return
	}

	if len(collection) == 0 {
		httpx.Error(w, http.StatusNotFound, "no data found", httpx.CodeNotFound)
		return
	}

	helper.SetResponse(w, http.StatusOK)
	_ = json.NewEncoder(w).Encode(collection)

}
