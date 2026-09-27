package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Sahilkumar121/workout_tracker/internal/helper"
	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
)

type ctxUserIDKey string

const authorizationkey = "Authorization"
const authorizationkeyID ctxUserIDKey = "userID"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get(authorizationkey)
		if authHeader == "" {
			httpx.Error(w, http.StatusUnauthorized, "missing authorized header", httpx.CodeUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			httpx.Error(w, http.StatusUnauthorized, "missing <type | token>", httpx.CodeUnauthorized)
			return
		}

		tokenString := parts[1]

		user_id, err := helper.VerifyAndDecodeToken(tokenString)
		if err != nil {
			httpx.Error(w, http.StatusUnauthorized, "invalid or expired token", httpx.CodeUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authorizationkeyID, user_id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthUserIDFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(authorizationkeyID).(string)
	return userID
}
