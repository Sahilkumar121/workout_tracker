package helper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sahilkumar121/workout_tracker/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SetResponse(w http.ResponseWriter, status int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}

func GetHashPassword(plainPassword string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func CheckHashPassword(hashPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))
}

func CreateAccessToken(userID string, secreteKey []byte) (string, error) {

	// 1. create payload
	claims := jwt.MapClaims{
		"sub":  userID,                               // Subject (User ID)
		"role": "auth",                               // Custom claim
		"iss":  "workout_tracker",                    // Issuer
		"iat":  time.Now().Unix(),                    // Issued at
		"exp":  time.Now().Add(time.Hour * 2).Unix(), // Expiration time (2 hours)
	}

	// hash the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the token with secrete key
	signedStr, err := token.SignedString(secreteKey)
	if err != nil {
		return "", err
	}

	return signedStr, nil
}

func VerifyAndDecodeToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(config.MustLoad().SecreteKey), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("subject claim missing or not a string")
	}

	return userID, nil
}
