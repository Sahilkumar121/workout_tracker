package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxRequestID string

const requestID = "X-Request-ID"

var requestIDKey ctxRequestID = "requestid"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := r.Header.Get(requestID)
		if id == "" {
			id = uuid.NewString()
		}

		w.Header().Add(requestID, id)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFromContext(ctx context.Context) string {
	return ctx.Value(requestIDKey).(string)
}
