package middleware

import (
	"merch-store/internal/common"
	"net/http"
	"github.com/golang-jwt/jwt"
)

const signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"

type Handler struct{}

func (h Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			common.WriteErrorResponse(w, http.StatusUnauthorized, "Токен не найден")
			return
		}

		tokenString := cookie.Value
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		})

		if err != nil || !token.Valid {
			common.WriteErrorResponse(w, http.StatusUnauthorized, "Неавторизован")
			return
		}

		next.ServeHTTP(w, r)
	})
}