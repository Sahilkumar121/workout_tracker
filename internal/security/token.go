package security

import (
	"errors"
	"time"

	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"github.com/golang-jwt/jwt/v5"
)

type customClaim struct {
	UserId string `json:"id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func CreateAccessToken(id string, secretKey []byte) (string, error) {

	claim := customClaim{
		UserId:    id,
		Role:      "client",
		Issuer:    "my-go-app",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(secretKey)
}

func ValidateAccessToken(tokenStr string, secretKey []byte) (*customClaim, error) {

	claim := &customClaim{}

	token, err := jwt.ParseWithClaims(tokenStr, claim, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, httpx.ErrWrongSigningMethod
		}

		return secretKey, nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer("my-go-app"),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, httpx.ErrTokenExpired
		}
		return nil, err
	}

	if !token.Valid {
		return nil, httpx.ErrInvalidToken
	}

	return claim, nil
}
