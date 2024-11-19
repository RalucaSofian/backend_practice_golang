package middlewares

import (
	"app/models"
	"app/services"
	"app/utils"
	"app/utils/access_utils"
	"context"
	"fmt"
	"net/http"
)

type ContextKey string

const (
	CURRENT_USER_KEY ContextKey = "current-user"
)

// Get the Current Auth User from the Context of an HTTP Request
func GetCurrentUser(req *http.Request) (*models.AuthUser, error) {
	contextUser := req.Context().Value(CURRENT_USER_KEY)

	currentUser, ok := contextUser.(*models.AuthUser)
	if !ok {
		fmt.Println("[middleware] Could not cast Current User")
		return nil, utils.NewApiError(utils.ErrorType_FormatError, "Could not cast Current User")
	}

	return currentUser, nil
}

// Auth Middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWr http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")

		authToken, err := access_utils.ExtractAuthToken(authHeader)
		if err != nil {
			http.Error(respWr, "Unauthorized", http.StatusUnauthorized)
			return
		}

		validToken := access_utils.ValidAccessToken(authToken)
		if !validToken {
			http.Error(respWr, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userEmail, err := access_utils.ExtractEmailFromToken(authToken)
		if err != nil {
			http.Error(respWr, "Unauthorized", http.StatusUnauthorized)
			return
		}

		currentUser, err := services.GetUserByEmail(userEmail)
		if err != nil {
			http.Error(respWr, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), CURRENT_USER_KEY, currentUser)

		next.ServeHTTP(respWr, req.WithContext(ctx))
	})
}
