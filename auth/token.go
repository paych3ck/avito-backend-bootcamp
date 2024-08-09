package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	jwt.StandardClaims
}

func GenerateJwtToken(email, userType string) (string, error) {
	expTime := time.Now().Add(72 * time.Hour)

	claims := &Claims{
		Email:    email,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtTokenString, err := jwtToken.SignedString(jwtSecretKey)

	if err != nil {
		return "", err
	}

	return jwtTokenString, nil
}

func ValidateJwtToken(jwtTokenStr string) (*Claims, error) {
	claims := &Claims{}

	jwtToken, err := jwt.ParseWithClaims(jwtTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, err
	}

	return claims, nil
}
