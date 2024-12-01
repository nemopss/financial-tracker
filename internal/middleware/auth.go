package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nemopss/financial-tracker/internal/response"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "Authorization header missing")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Error(w, http.StatusUnauthorized, "Invalid Authorization header format")
				return
			}
			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, http.ErrAbortHandler
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				log.Printf("Invalid token: %v", err)
				response.Error(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || claims["user_id"] == nil {
				response.Error(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			userID, ok := claims["user_id"].(float64)
			if !ok {
				response.Error(w, http.StatusUnauthorized, "Invalid user ID")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, int(userID))
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)

		})
	}
}
