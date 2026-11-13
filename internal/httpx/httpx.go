package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type customError string

var (
	ErrInternalError customError = "internal server error"
	ErrNotFound      customError = "not found"
	ErrBadRequest    customError = "bad request"
)

type errorStruct struct {
	error customError
	msg   string
}

type errorEnvelope struct {
	error errorStruct
}

func HandlerError(w http.ResponseWriter, status int, msg string, error customError) {

	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(errorEnvelope{
		error: errorStruct{
			error: error,
			msg:   msg,
		},
	}); err != nil {
		slog.Error("httpx HandlerError encode failed", "err", err)
		return
	}
}
