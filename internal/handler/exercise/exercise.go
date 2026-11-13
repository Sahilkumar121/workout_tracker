package exercise

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	Logger  *slog.Logger
	Service Service
}

type Service interface {
	CreateExercise(
		ctx context.Context,
		exerciseBodyReq PostRequest,
	) (Exercise, error)

	FindAllExercise(
		ctx context.Context,
	) ([]Exercise, error)

	SearchExerciseByID(
		ctx context.Context,
		id string,
	) (Exercise, error)
}

func NewExerciseHandler(
	logger *slog.Logger,
	service Service,
) *Handler {
	return &Handler{
		Logger:  logger,
		Service: service,
	}
}

func (h *Handler) PostExercise(
	w http.ResponseWriter, r *http.Request,
) {

	var req PostRequest

	//* decode the request data and return bad request if error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("failed post exercise decode", "error", err)
		httpx.HandlerError(
			w,
			http.StatusInternalServerError,
			"server error, try again later",
			httpx.ErrInternalError,
		)
		return
	}

	response, err := h.Service.CreateExercise(r.Context(), req)
	if err != nil {
		h.Logger.Error("create exercise failed", "err", err)

		if errors.Is(err, httpx.ErrBadRequest) {
			httpx.HandlerError(w, http.StatusBadRequest, "please provide correct field and values", httpx.ErrBadRequest)
			return
		}

		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			httpx.HandlerError(w, http.StatusConflict, "exercise already exist", httpx.ErrConflict)
			return
		}
		httpx.HandlerError(w, http.StatusInternalServerError, "please try again later", httpx.ErrInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func (h *Handler) GetExercises(w http.ResponseWriter, r *http.Request) {
	d, err := h.Service.FindAllExercise(r.Context())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.Logger.Error("failed get exercises", "err", err)
			httpx.HandlerError(w, http.StatusNotFound, "not exercise exists in database", httpx.ErrNotFound)
			return
		}
		h.Logger.Error("failed get exercises", "err", err)
		httpx.HandlerError(w, http.StatusInternalServerError, "server side error, try again after sometime", httpx.ErrInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		h.Logger.Error("json encode error", "err", err)
		httpx.HandlerError(w, http.StatusInternalServerError, "server side error, try again after sometime", httpx.ErrInternalError)
		return
	}
}

/*

Todo: get exercise by id (Param query)
* id
* name
* description
* muscle_group
* created_at

*/

func (h *Handler) GetExerciseById(w http.ResponseWriter, r *http.Request) {

	reqId := r.PathValue("id")

	d, err := h.Service.SearchExerciseByID(r.Context(), reqId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.Logger.Error("failed get exercise by id", "err", err)
			httpx.HandlerError(w, http.StatusNotFound, "no exercise found", httpx.ErrNotFound)
			return
		}
		h.Logger.Error("failed get exercise by id", "err", err)
		httpx.HandlerError(w, http.StatusInternalServerError, "server error, try again later", httpx.ErrInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		h.Logger.Error("failed decode json", "err", err)
		httpx.HandlerError(w, http.StatusInternalServerError, "server error, try again later", httpx.ErrInternalError)
		return
	}
}

/*

Todo: patch exercise by id (Param query)

* name - optional
* description - optional
* muscle_group - optional

*/

/*

Todo: delete exercise by id (Param query)

? return json deleted successfully (status code 200)

*/
