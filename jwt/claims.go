package jwt

import (
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type CustomClaims struct {
	TokenType string `json:"token_type"`
	UserID    int64  `json:"user_id"`
	jwt5.RegisteredClaims
}

func GenerateAccessTokenJWT(userID int64, secret string, tokenValidity time.Duration, logger *logrus.Logger) (string, error) {
	token := jwt5.New(jwt5.SigningMethodHS256)

	claims := &CustomClaims{
		TokenType: "access",
		UserID:    userID,
		RegisteredClaims: jwt5.RegisteredClaims{
			ExpiresAt: jwt5.NewNumericDate(time.Now().Add(tokenValidity)),
			IssuedAt:  jwt5.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}
	token.Claims = claims

	// Generate encoded token and send it as a response.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		logger.Errorf("Error signing the token: %v", err)
		return "", err
	}
	return t, nil
}
