package jwt

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// TokenIssuer interface to generate and reissue tokens.
type TokenIssuer interface {
	NewToken(db *gorm.DB, userID int64, refreshTokenValidity time.Duration) (string, error)
	ReissueToken(db *gorm.DB, oldToken string, refreshTokenValidity time.Duration) (string, int64, error)
}

func NewPairToken(db *gorm.DB, secret string, userID int64, issuer TokenIssuer, tokenValidity time.Duration, refreshTokenValidity time.Duration, logger *logrus.Logger) (*Token, error) {
	var err error
	accessToken, err := GenerateAccessTokenJWT(userID, secret, tokenValidity, logger)
	if err != nil {
		return nil, err
	}
	refreshToken, err := issuer.NewToken(db, userID, refreshTokenValidity)
	if err != nil {
		return nil, err
	}
	return &Token{accessToken, refreshToken}, nil
}

func RefreshPairToken(db *gorm.DB, secret string, oldToken string, issuer TokenIssuer, tokenValidity time.Duration, refreshTokenValidity time.Duration, logger *logrus.Logger) (*Token, error) {
	var err error
	refreshToken, userID, err := issuer.ReissueToken(db, oldToken, refreshTokenValidity)
	if err != nil {
		return nil, err
	}

	accessToken, err := GenerateAccessTokenJWT(userID, secret, tokenValidity, logger)
	if err != nil {
		return nil, err
	}
	return &Token{accessToken, refreshToken}, nil
}
