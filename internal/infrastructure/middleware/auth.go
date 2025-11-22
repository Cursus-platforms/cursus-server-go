package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Cursus-platforms/cursus-server-go/internal/lib/response"
)

type ContextKey string

const UserIDKey ContextKey = "userID"
const UserRoleKey ContextKey = "userRole"

var jwtSecretKey = []byte("YOUR_SUPER_SECRET_KEY")

func RequireRole(allowedRoles ...string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.RespondWithError(w, http.StatusUnauthorized, "Missing Authorization header", nil)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				response.RespondWithError(w, http.StatusUnauthorized, "Invalid token format", nil)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtSecretKey, nil
			})

			if err != nil || !token.Valid {
				response.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token", nil)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				response.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims", nil)
				return
			}

			userIDStr, okSub := claims["sub"].(string)
			userRole, okRole := claims["role"].(string)
			if !okSub || !okRole {
				response.RespondWithError(w, http.StatusUnauthorized, "Token is missing user ID or role", nil)
				return
			}

			isAllowed := false
			for _, role := range allowedRoles {
				if userRole == role {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				log.Printf("ACCESS DENIED: User %s (Role: %s) attempted unauthorized access", userIDStr, userRole)
				response.RespondWithError(w, http.StatusForbidden, "Access forbidden: Insufficient role", nil)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userIDStr)
			ctx = context.WithValue(ctx, UserRoleKey, userRole)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
