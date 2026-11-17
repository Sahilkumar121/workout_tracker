package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"github.com/Sahilkumar121/workout_tracker/internal/security"
)

const userCtxID = "userID"

func AuthenticateUser(secretKey []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		fields := strings.Fields(authHeader)

		if len(fields) < 2 || strings.EqualFold(fields[0], "Bearer") {
			httpx.HandleResponseError(w, httpx.ErrUnauthorized)
			return
		}

		token := fields[1]

		claim, err := security.ValidateAccessToken(token, secretKey)
		if err != nil {
			httpx.HandleResponseError(w, httpx.ErrUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userCtxID, claim.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) string {
	return ctx.Value(userCtxID).(string)
}
