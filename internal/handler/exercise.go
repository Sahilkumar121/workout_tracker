package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Sahilkumar121/workout_tracker/internal/helper"
	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"github.com/Sahilkumar121/workout_tracker/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExerciseHandler struct {
	DB       *pgxpool.Pool
	Logger   *slog.Logger
	Validate *validator.Validate
}

func NewExerciseHandler(db *pgxpool.Pool, logger *slog.Logger) *ExerciseHandler {
	return &ExerciseHandler{
		DB:       db,
		Logger:   logger,
		Validate: validator.New(),
	}
}

func (h *ExerciseHandler) PostExercises(w http.ResponseWriter, r *http.Request) {

	requestID := middleware.RequestIDFromContext(r.Context())

	var requestData ExerciseRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		h.Logger.Error("failed decode exercise", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusInternalServerError, "data decode error", httpx.CodeInternalServer)
		return
	}

	err = h.Validate.Struct(requestData)
	if err != nil {
		h.Logger.Error("failed validate, bad request", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusBadRequest, "please provide information properly", httpx.CodeInvalidInput)
		return
	}

	query := `insert into exercises (exercise_name, target_muscle) values ($1,$2)`
	args := []any{requestData.ExerciseName, requestData.TargetMuscle}

	_, err = h.DB.Exec(r.Context(), query, args...)
	if err != nil {

	}
}

func (h *ExerciseHandler) GetAllExercise(w http.ResponseWriter, r *http.Request) {

	query := `select id, exercise_name, target_muscle, created_at from exercises`

	rows, err := h.DB.Query(r.Context(), query)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalServer)
		return
	}

	collections, err := pgx.CollectRows(rows, pgx.RowToStructByName[ExerciseResponse])

	if len(collections) == 0 {
		httpx.Error(w, http.StatusNotFound, "no data found", httpx.CodeNotFound)
		return
	}

	helper.SetResponse(w, http.StatusOK)
	_ = json.NewEncoder(w).Encode(collections)
}

// * this method is used to get the exercise detail using it's id
func (h *ExerciseHandler) GetExerciseByID(w http.ResponseWriter, r *http.Request) {

	requestID := middleware.RequestIDFromContext(r.Context())
	userID := middleware.AuthUserIDFromContext(r.Context())

	exerciseID := r.PathValue("id")

	query := `select id, exercise_name, target_muscle, created_at from exercises where id=%s`

	rows := h.DB.QueryRow(r.Context(), query, exerciseID)

	var responseData ExerciseResponse
	if err := rows.Scan(&responseData); err != nil {

		h.Logger.Error("failed to scan rows", "err", err, "requestID", requestID, "userID", userID) //! this is alert
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalServer)
		return
	}

}
