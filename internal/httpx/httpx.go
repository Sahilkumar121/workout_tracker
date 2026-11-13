package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type customError string

func (e customError) Error() string {
	return string(e)
}

var (
	ErrInternalError    customError = "internal server error"
	ErrNotFound         customError = "not found"
	ErrBadRequest       customError = "bad request"
	ErrConflict         customError = "conflict data"
	ErrNoFieldsToUpdate customError = "no filed to update"
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
