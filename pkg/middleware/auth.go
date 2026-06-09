package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

const ContextUserID = "user_id"

type AuthMiddleware struct {
	jwtSecret string
	Logger *zap.Logger
}

func NewAuthMiddleware(jwtSecret string, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		Logger: logger,
	}
}

func (m *AuthMiddleware) IsAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		tokenString, _ := strings.CutPrefix(header, "Bearer ")

		if tokenString == "" {
			m.Logger.Warn("missing token")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return []byte(m.jwtSecret), nil
		})
		if err != nil {
			m.Logger.Warn("missing token")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			m.Logger.Warn("invalid token")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return 
		}

		ctx := context.WithValue(r.Context(), ContextUserID, claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
