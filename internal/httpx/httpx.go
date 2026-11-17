package httpx

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type customError string

func (e customError) Error() string {
	return string(e)
}

var (
	// hanlder
	ErrInternalError    customError = "internal server error"
	ErrNotFound         customError = "not found"
	ErrBadRequest       customError = "bad request"
	ErrConflict         customError = "conflict data"
	ErrNoFieldsToUpdate customError = "no filed to update"

	// hash password
	ErrPasswordLen  customError = "incorrect password len"
	ErrArgonVersion customError = "argon version wrong"

	// access token
	ErrUnauthorized       customError = "unauthorized"
	ErrWrongSigningMethod customError = "wrong signed method"
	ErrInvalidToken       customError = "invalid token"
	ErrTokenExpired       customError = "token expired"
	ErrInvalidPasswordHash customError = "invalid password hash"
)

type errorStruct struct {
	Type customError `json:"type"`
	Msg  string      `json:"msg"`
}

type errorEnvelope struct {
	Error errorStruct `json:"error"`
}

func HandlerError(w http.ResponseWriter, status int, msg string, errorType customError) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(errorEnvelope{
		Error: errorStruct{
			Type: errorType,
			Msg:  msg,
		},
	}); err != nil {
		slog.Error("httpx HandlerError encode failed", "err", err)
	}
}

func HandleResponseError(w http.ResponseWriter, err error) {
	var pgErr *pgconn.PgError

	switch {
	case errors.Is(err, ErrBadRequest):
		HandlerError(
			w,
			http.StatusBadRequest,
			"please provide proper info",
			ErrBadRequest,
		)

	case errors.Is(err, ErrNoFieldsToUpdate):
		HandlerError(
			w,
			http.StatusBadRequest,
			"please provided at least one field",
			ErrNoFieldsToUpdate,
		)

	case errors.Is(err, ErrUnauthorized):
		HandlerError(
			w,
			http.StatusUnauthorized,
			"invlide credential",
			ErrUnauthorized,
		)

	case errors.Is(err, ErrWrongSigningMethod):
		HandlerError(
			w,
			http.StatusUnauthorized,
			"token method is wrong",
			ErrWrongSigningMethod,
		)

	case errors.Is(err, ErrInvalidToken):
		HandlerError(
			w,
			http.StatusUnauthorized,
			"token is wrong",
			ErrInvalidToken,
		)

	case errors.Is(err, ErrTokenExpired):
		HandlerError(
			w,
			http.StatusUnauthorized,
			"token has been expired, please renew token",
			ErrTokenExpired,
		)

	case errors.Is(err, pgx.ErrNoRows):
		HandlerError(
			w,
			http.StatusNotFound,
			"no data to be found",
			ErrNotFound,
		)

	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		HandlerError(
			w,
			http.StatusConflict,
			"requested data already exist",
			ErrConflict,
		)

	default:
		HandlerError(
			w,
			http.StatusInternalServerError,
			"server error, try again",
			ErrInternalError,
		)
	}
}
