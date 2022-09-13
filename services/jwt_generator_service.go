package services

import (
	"Ferrum/data"
	"github.com/golang-jwt/jwt/v4"
)

type JwtGenerator struct {
}

func (generator *JwtGenerator) GenerateAccessToken(tokenData *data.AccessTokenData) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)
	// todo(UMV): where we get key ... ???
	signedToken, err := token.SignedString([]byte("secureSecretText"))
	if err != nil {
		//todo(UMV): think what to do on Error
	}
	return signedToken
}

func (generator *JwtGenerator) GenerateRefreshToken(tokenData *data.TokenRefreshData) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)
	// todo(UMV): where we get key ... ???
	signedToken, err := token.SignedString([]byte("secureSecretText"))
	if err != nil {
		//todo(UMV): think what to do on Error
	}
	return signedToken
}
