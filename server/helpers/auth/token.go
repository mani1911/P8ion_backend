package helper

import (
	"errors"
	"fmt"
	"p8ion/config"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func GenerateToken(userID uint) (string, error) {
	Auth := config.GetConfig().Auth
	tokenLifespan, err := strconv.Atoi(Auth.TokenHourLifeSpan)

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(Auth.JWTSecret))

}

func ValidateToken(authHeader string) (uint, error) {
	Auth := config.GetConfig().Auth
	if Auth.JWTSecret == "" {
		return 0, errors.New("JWT_SECRET not set")
	}

	if authHeader == "" {
		return 0, errors.New("Missing Authorization Header")
	}

	if len(authHeader) < 7 || !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, errors.New("Invalid Authorization Header")
	}

	tokenString := authHeader[7:]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(Auth.JWTSecret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("Invalid JWT")
	}

	userID, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["userID"]), 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(userID), nil
}
