package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/Sahilkumar121/workout_tracker/internal/helper"
)

type codeErrorString string

const (
	CodeInternalServer codeErrorString = "internal_server"
	CodeInvalidInput   codeErrorString = "bad_request"
	CodeUnauthorized   codeErrorString = "unauthorized"
	CodeConflict       codeErrorString = "data_conflict"
	CodeNotFound       codeErrorString = "data_not_found"
)

type errorEnvelop struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Message string          `json:"message"`
	Code    codeErrorString `json:"code"`
}

func Error(w http.ResponseWriter, status int, message string, code codeErrorString) {

	helper.SetResponse(w, status)
	_ = json.NewEncoder(w).Encode(errorEnvelop{Error: errorPayload{Message: message, Code: code}})
}
