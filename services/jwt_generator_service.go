package services

import (
	"Ferrum/data"
	"github.com/golang-jwt/jwt/v4"
)

type JwtGenerator struct {
}

func (generator *JwtGenerator) GenerateAccessToken(tokenData *data.AccessTokenData) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)
	// todo(UMV): should we sign token or not ???
	return token.Raw
}

func (generator *JwtGenerator) GenerateRefreshToken(tokenData *data.TokenRefreshData) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)
	// todo(UMV): should we sign token or not ???
	return token.Raw
}
