package access_utils

import (
	"app/utils"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: env var / constants
var secretKey = []byte("SECRET_KEY")

// Create an Access Token, based on some User Data
func CreateAccessToken(dataToEncode any) (string, error) {
	tokenClaims := jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * 4).Unix(),
		"data": dataToEncode,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	tokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		fmt.Println("[utils]", err.Error())
		apiErr := utils.NewApiError(utils.ErrorType_JWTError, err.Error())
		return "", apiErr
	}
	return tokenString, nil
}

// Validate an Access Token
func ValidAccessToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		fmt.Println("[utils]", err.Error())
		return false
	}
	if !token.Valid {
		return false
	}

	return true
}

// Extract the Auth User Email from an Access Token
func ExtractEmailFromToken(authToken string) (string, error) {
	token, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		fmt.Println("[utils]", err.Error())
		return "", err
	}
	if !token.Valid {
		return "", err
	}

	userEmail := token.Claims.(jwt.MapClaims)["data"]

	return userEmail.(string), nil
}

// Extract the Token from the Authorization Header of a HTTP Request
func ExtractAuthToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", utils.NewApiError(utils.ErrorType_JWTError, "Invalid Token or Auth Type")
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) < 2 {
		return "", utils.NewApiError(utils.ErrorType_JWTError, "Invalid Token or Auth Type")
	}

	authType := authParts[0]
	authToken := authParts[1]
	if authType != "Bearer" || authToken == "" {
		return "", utils.NewApiError(utils.ErrorType_JWTError, "Invalid Token or Auth Type")
	}

	return authToken, nil
}
