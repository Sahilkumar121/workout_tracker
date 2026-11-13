package httpx

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrInternalError = errors.New("internal server error")
	ErrNotFound      = errors.New("not found")
	ErrBadRequest    = errors.New("bad request")
	ErrConflict      = errors.New("conflict data")
)

type errorStruct struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

type errorEnvelope struct {
	Error errorStruct `json:"error"`
}

func HandlerError(w http.ResponseWriter, status int, msg string, errorType error) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(errorEnvelope{
		Error: errorStruct{
			Type: errorType.Error(),
			Msg:  msg,
		},
	}); err != nil {
		slog.Error("httpx HandlerError encode failed", "err", err)

	}
}
