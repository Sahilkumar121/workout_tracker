package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const requestIdKey = "X-Request-ID"

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestId := r.Header.Get(requestIdKey)
		if strings.TrimSpace(requestId) == "" {
			requestId = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), requestIdKey, requestId)

		w.Header().Add(requestIdKey, requestId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestId(ctx context.Context) string {
	return ctx.Value(requestIdKey).(string)
}
