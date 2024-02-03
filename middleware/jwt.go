package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/aredcomet/go-gears/jwt"
	"github.com/aredcomet/go-gears/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"

	jwt5 "github.com/golang-jwt/jwt/v5"
)

func extractBearerToken(r *http.Request, headerName string, authScheme string) (string, bool) {
	authHeader := r.Header.Get(headerName)
	if authHeader == "" {
		return "", false
	}
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || parts[0] != authScheme {
		return "", false
	}
	return parts[1], true
}

func customKeyFunc(secret string) jwt5.Keyfunc {
	return func(token *jwt5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt5.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}
}

func JWTMiddleware(claimKey string, secret string, headerName string, authScheme string, tokenType string, logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, ok := extractBearerToken(r, headerName, authScheme)
			if !ok {
				utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header is missing or invalid", logger)
				return
			}

			token, err := jwt5.ParseWithClaims(tokenString, &jwt.CustomClaims{}, customKeyFunc(secret))
			if err != nil {

				if errors.Is(err, jwt5.ErrTokenExpired) {
					utils.RespondWithError(w, http.StatusUnauthorized, "Your token is expired", logger)
					return
				}
				utils.RespondWithError(w, http.StatusBadRequest, err.Error(), logger)
				return
			}

			if !token.Valid {
				utils.RespondWithError(w, http.StatusUnauthorized, "Token validation failed", logger)
				return
			}

			claims, ok := token.Claims.(*jwt.CustomClaims)
			if !ok || claims.TokenType != tokenType {
				utils.RespondWithError(w, http.StatusUnauthorized, "Unable to extract claims or token type is wrong", logger)
				return
			}
			ctx := context.WithValue(r.Context(), claimKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
