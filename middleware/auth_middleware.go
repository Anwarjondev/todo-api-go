package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/Anwarjondev/todo-api-go/handlers"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtkey []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error with loading .env file")
	}
	jwtkey = []byte(os.Getenv("JWT_KEY"))
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtkey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		if claims.UserId == 0 {
			http.Error(w, "Unathorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", int(claims.UserId))
		ctx = context.WithValue(ctx, "role", string(claims.Role))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}