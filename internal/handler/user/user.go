package user

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"github.com/Sahilkumar121/workout_tracker/internal/middleware"
	"github.com/Sahilkumar121/workout_tracker/internal/security"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Logger    *slog.Logger
	validate  *validator.Validate
	service   UserService
	secretKey []byte
}

func NewUserHandler(logger *slog.Logger, service UserService, secretKey []byte) *Handler {
	return &Handler{
		Logger:    logger,
		validate:  validator.New(validator.WithRequiredStructEnabled()),
		service:   service,
		secretKey: secretKey,
	}
}

type UserService interface {
	CheckValidUserService(
		ctx context.Context,
		email string,
	) (*UserCredentials, error)

	RegisterUserService(
		ctx context.Context,
		req *RegisterRequest,
	) (*UserResponse, error)

	FindUserByID(
		ctx context.Context,
		id string,
	) (*UserResponse, error)
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	requestID := middleware.GetRequestId(r.Context())

	var registerReq RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		h.Logger.Error("failed decode json request", "err", err, "requestID", requestID)
		httpx.HandleResponseError(w, err)
		return
	}

	// validate the struct
	if err := h.validate.Struct(registerReq); err != nil {
		h.Logger.Error("failed vaidate register request", "err", err, "requestID", requestID)
		httpx.HandleResponseError(w, httpx.ErrBadRequest)
		return
	}

	d, err := h.service.RegisterUserService(r.Context(), &registerReq)
	if err != nil {
		h.Logger.Error("failed register user", "err", err, "requestID", requestID)
		httpx.HandleResponseError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(d); err != nil {
		h.Logger.Error("decode error", "err", err, "requestID", requestID)
	}
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {

	requestID := middleware.GetRequestId(r.Context())

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		h.Logger.Error("failed decode request", "err", err, "requestID", requestID)
		httpx.HandleResponseError(w, err)
		return
	}

	if err := h.validate.Struct(loginReq); err != nil {
		h.Logger.Error("failed struct", "err", err, "requestID", requestID)
		httpx.HandleResponseError(w, httpx.ErrBadRequest)
		return
	}

	d, err := h.service.CheckValidUserService(r.Context(), loginReq.Email)
	if err != nil {
		h.Logger.Error("no user found", "err", err, "requestID", requestID)
		httpx.HandleResponseError(w, err)
		return
	}

	isCorrect, err := security.VerifyPassword(loginReq.Password, d.Password)
	if err != nil {
		h.Logger.Error("login user", "err", err, "requestID", requestID)
		httpx.HandleResponseError(w, httpx.ErrUnauthorized)
		return
	}

	if !isCorrect {
		h.Logger.Error("wrong password enter", "requestID", requestID)
		httpx.HandleResponseError(w, httpx.ErrUnauthorized)
		return
	}

	t, err := security.CreateAccessToken(d.ID.String(), h.secretKey)
	if err != nil {
		h.Logger.Error("failed access token", "err", err, "requestID", requestID)
		httpx.HandleResponseError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"msg":        "login successfull",
		"token":      t,
		"token-type": "Bearer",
	})

}

// create user me
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {

	requestID := middleware.GetRequestId(r.Context())
	userID := middleware.GetUserID(r.Context())

	d, err := h.service.FindUserByID(r.Context(), userID)
	if err != nil {
		h.Logger.Error("failed find user by id", "err", err, "requestID", requestID)
		httpx.HandleResponseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(d); err != nil {
		h.Logger.Error("failed decode", "err", err, "requestID", requestID)
	}
}
