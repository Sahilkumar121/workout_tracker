package exercise

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"github.com/Sahilkumar121/workout_tracker/internal/middleware"
	"github.com/go-playground/validator/v10"
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

	ValidateExerciseId(
		id string,
	) error

	DeleteExerciseById(
		ctx context.Context,
		id string,
	) (pgconn.CommandTag, error)
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

	requestID := middleware.GetRequestId(r.Context())

	var req PostRequest

	//* decode the request data and return bad request if error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("failed post exercise decode", "error", err, "requestID", requestID)
		httpx.HandleResponseError(w, err)
		return
	}

	// validate the PostRequest
	if err := h.validate.Struct(req); err != nil {
		h.Logger.Error("create exercise failed", "err", err)
		httpx.HandleResponseError(w, err)
		return
	}

	response, err := h.Service.CreateExercise(r.Context(), req)
	if err != nil {
		h.Logger.Error("create exercise failed", "err", err)
		httpx.HandleResponseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.Error("failed encoded", "err", err)
	}
}

// GetExercises return all available exercise
func (h *Handler) GetExercises(
	w http.ResponseWriter, r *http.Request,
) {

	d, err := h.Service.FindAllExercise(r.Context())
	if err != nil {
		h.Logger.Error("failed get exercises", "err", err)
		httpx.HandleResponseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(d); err != nil {
		h.Logger.Error("json encode error", "err", err)
	}
}

// GetExerciseById return the exercise with matching id
func (h *Handler) GetExerciseById(w http.ResponseWriter, r *http.Request) {

	reqId := r.PathValue("id")

	if err := h.Service.ValidateExerciseId(reqId); err != nil {
		h.Logger.Error("get exercise by id", "err", err)
		httpx.HandleResponseError(w, err)
		return
	}

	d, err := h.Service.SearchExerciseByID(r.Context(), reqId)
	if err != nil {
		h.Logger.Error("failed get exercise by id", "err", err)
		httpx.HandleResponseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(d); err != nil {
		h.Logger.Error("failed encode json", "err", err)
	}
}

// PatchExerciseById update the matching exercise data
func (h *Handler) PatchExerciseById(w http.ResponseWriter, r *http.Request) {

	// path query id
	exerciseId := r.PathValue("id")

	if err := h.Service.ValidateExerciseId(exerciseId); err != nil {
		h.Logger.Error("patch exercise by id", "err", err)
		httpx.HandleResponseError(w, err)
		return
	}

	// exercise struct for update
	var patchReq PatchRequest
	// decode from request
	if err := json.NewDecoder(r.Body).Decode(&patchReq); err != nil {
		h.Logger.Error("failed decode patchReq", "err", err)
		httpx.HandleResponseError(w, err)
		return
	}

	d, err := h.Service.UpdateExerciseById(r.Context(), exerciseId, patchReq)
	if err != nil {
		h.Logger.Error("patch exercise by id", "err", err)
		httpx.HandleResponseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	h.Logger.Info("updated exercise successfully", "id", exerciseId)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		h.Logger.Error("failed encode json", "err", err)
	}
}

func (h *Handler) DeleteExerciseById(
	w http.ResponseWriter,
	r *http.Request,
) {

	exerciseId := r.PathValue("id")

	if err := h.Service.ValidateExerciseId(exerciseId); err != nil {
		h.Logger.Error("patch exercise by id", "err", err)
		httpx.HandleResponseError(w, err)
		return
	}

	c, err := h.Service.DeleteExerciseById(r.Context(), exerciseId)
	if err != nil {
		h.Logger.Error("delete exercise by id", "err", err)
		httpx.HandleResponseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	h.Logger.Info("delete exercise", "id", exerciseId, "rows_affected", c.RowsAffected())
	if err := json.NewEncoder(w).Encode(map[string]any{
		"msg":           "successfully delete exercise",
		"rows_affected": c.RowsAffected(),
	}); err != nil {
		h.Logger.Error("encode json error", "err", err)
	}
}
