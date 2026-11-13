package exercise

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	Logger   *slog.Logger
	Service  Service
	validate *validator.Validate
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

	UpdateExerciseById(
		ctx context.Context,
		id string,
		patchReq PatchRequest,
	) (Exercise, error)
}

func NewExerciseHandler(
	logger *slog.Logger,
	service Service,
) *Handler {
	return &Handler{
		Logger:   logger,
		Service:  service,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// PostExercise - HTTP POST req for exercise
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

	// validate the PostRequest
	if err := h.validate.Struct(req); err != nil {
		h.Logger.Error("create exercise failed", "err", err)
		httpx.HandlerError(
			w,
			http.StatusBadRequest,
			"please provide correct field and values",
			httpx.ErrBadRequest,
		)
		return
	}

	response, err := h.Service.CreateExercise(r.Context(), req)
	if err != nil {
		h.Logger.Error("create exercise failed", "err", err)

		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			httpx.HandlerError(
				w,
				http.StatusConflict,
				"exercise already exist",
				httpx.ErrConflict,
			)
			return
		}
		httpx.HandlerError(
			w,
			http.StatusInternalServerError,
			"please try again later",
			httpx.ErrInternalError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func (h *Handler) GetExercises(
	w http.ResponseWriter, r *http.Request,
) {

	d, err := h.Service.FindAllExercise(r.Context())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.Logger.Error("failed get exercises", "err", err)
			httpx.HandlerError(
				w,
				http.StatusNotFound,
				"not exercise exists in database",
				httpx.ErrNotFound,
			)
			return
		}
		h.Logger.Error("failed get exercises", "err", err)
		httpx.HandlerError(
			w,
			http.StatusInternalServerError,
			"server side error, try again after sometime",
			httpx.ErrInternalError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		h.Logger.Error("json encode error", "err", err)
		httpx.HandlerError(
			w,
			http.StatusInternalServerError,
			"server side error, try again after sometime",
			httpx.ErrInternalError,
		)
		return
	}
}

func (h *Handler) GetExerciseById(w http.ResponseWriter, r *http.Request) {

	reqId := r.PathValue("id")

	if _, err := uuid.Parse(reqId); err != nil {
		h.Logger.Error("get exercise by id", "err", err)
		httpx.HandlerError(w, http.StatusBadRequest, "please provide correct id", httpx.ErrBadRequest)
		return
	}

	d, err := h.Service.SearchExerciseByID(r.Context(), reqId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.Logger.Error("failed get exercise by id", "err", err)
			httpx.HandlerError(
				w,
				http.StatusNotFound,
				"no exercise found",
				httpx.ErrNotFound,
			)
			return
		}
		h.Logger.Error("failed get exercise by id", "err", err)
		httpx.HandlerError(
			w,
			http.StatusInternalServerError,
			"server error, try again later",
			httpx.ErrInternalError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		h.Logger.Error("failed decode json", "err", err)
		httpx.HandlerError(
			w,
			http.StatusInternalServerError,
			"server error, try again later",
			httpx.ErrInternalError,
		)
		return
	}
}

func (h *Handler) PatchExerciseById(w http.ResponseWriter, r *http.Request) {

	// path query id
	exerciseId := r.PathValue("id")

	// check id path query id was uuid format else return 400
	if _, err := uuid.Parse(exerciseId); err != nil || strings.TrimSpace(exerciseId) == "" {
		h.Logger.Error("parse exerciseId failed", "err", err)
		httpx.HandlerError(
			w,
			http.StatusBadRequest,
			"wrong exercise id",
			httpx.ErrBadRequest,
		)
		return
	}

	// check if the exercise exist or not
	_, err := h.Service.SearchExerciseByID(r.Context(), exerciseId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.HandlerError(
				w,
				http.StatusNotFound,
				"no exercise found for given id",
				httpx.ErrNotFound,
			)
			return
		}
		h.Logger.Error("patch exercise by id", "err", err)
		httpx.HandlerError(
			w,
			http.StatusInternalServerError,
			"server error, try again later",
			httpx.ErrInternalError,
		)
		return
	}

	// exercise struct for update
	var patchReq PatchRequest
	// decode from request
	if err := json.NewDecoder(r.Body).Decode(&patchReq); err != nil {
		h.Logger.Error("failed decode patchReq", "err", err)
		httpx.HandlerError(
			w,
			http.StatusInternalServerError,
			"server error, try again later",
			httpx.ErrInternalError,
		)
		return
	}

	d, err := h.Service.UpdateExerciseById(r.Context(), exerciseId, patchReq)
	if err != nil {
		h.Logger.Error("patch exercise by id", "err", err)
		if errors.Is(err, httpx.ErrNoFieldsToUpdate) {
			httpx.HandlerError(
				w,
				http.StatusBadRequest,
				"please provide at least one field",
				httpx.ErrNoFieldsToUpdate,
			)
			return
		}

		if errors.Is(err, pgx.ErrNoRows) {
			httpx.HandlerError(
				w,
				http.StatusNotFound,
				"no exercise found",
				httpx.ErrNoFieldsToUpdate,
			)
			return
		}
		httpx.HandlerError(
			w,
			http.StatusInternalServerError,
			"server error, try again later",
			httpx.ErrInternalError,
		)
		return
	}

	w.Header().Set("Content-Type", "application")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		return
	}
}

/*

Todo: delete exercise by id (Param query)

? return json deleted successfully (status code 200)

*/
