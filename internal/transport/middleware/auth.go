package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/matveevaolga/request-managing-app/internal/service"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	UserRoleKey contextKey = "userRole"
)

func CheckAdmin(next http.HandlerFunc, authService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-API-TOKEN")
		if token == "" {
			http.Error(w, `{"error": "missing X-API-TOKEN header"}`, http.StatusUnauthorized)
			return
		}

		claims, err := authService.ValidateToken(token)
		if err != nil {
			slog.Error("token validation failed", "error", err)
			http.Error(w, `{"error": "invalid token"}`, http.StatusUnauthorized)
			return
		}

		if claims.Role != "ADMIN" {
			http.Error(w, `{"error": "admin role required"}`, http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

}
