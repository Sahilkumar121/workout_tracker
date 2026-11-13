package exercise

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
)

type Handler struct {
	Logger  *slog.Logger
	Service Service
}

type Service interface {
	CreateExercise(
		ctx context.Context,
		exerciseBodyReq PostRequest,
	) (*Exercise, error)
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

/*
Todo: post new custom exercise
* name
* description
* muscle_group

*/

func (h *Handler) PostExercise(
	w http.ResponseWriter, r *http.Request,
) {

	var req PostRequest

	//* decode the request data and return bad request if error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("bad body request", "error", err)
		httpx.HandlerError(
			w,
			http.StatusBadRequest,
			"please provide all information",
			httpx.ErrBadRequest,
		)
		return
	}

	response, err := h.Service.CreateExercise(r.Context(), req)
	if err != nil {
		h.Logger.Error("create exercise failed: %v\n", err)
		httpx.HandlerError(w, http.StatusInternalServerError, "please try again later", httpx.ErrInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

/*

Todo: get all exercise
* id
* name
* description
* muscle_group
* created_at

*/

/*

Todo: get exercise by id (Param query)
* id
* name
* description
* muscle_group
* created_at

*/

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
