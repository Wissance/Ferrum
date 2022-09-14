package services

import (
	"Ferrum/data"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

type JwtGenerator struct {
}

func (generator *JwtGenerator) GenerateAccessToken(tokenData *data.AccessTokenData) string {
	token := jwt.New(jwt.SigningMethodHS256)
	key := []byte("secureSecretText")
	// todo(UMV): where we get key ... ???
	// signed token contains embedded type because we don't actually know type of User, therefore we do it like jwt do but use RawStr
	signedToken, err := generator.makeSignedToken(token, tokenData, key)
	//token.SignedString([]byte("secureSecretText"))
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

func (generator *JwtGenerator) makeSignedToken(token *jwt.Token, tokenData *data.AccessTokenData, signKey interface{}) (string, error) {
	var err error
	var sig string
	var jsonValue []byte

	if jsonValue, err = json.Marshal(token.Header); err != nil {
		return "", err
	}
	header := base64.RawURLEncoding.EncodeToString(jsonValue)

	claim := base64.RawURLEncoding.EncodeToString([]byte(tokenData.ResultJsonStr))

	unsignedToken := strings.Join([]string{header, claim}, ".")
	if sig, err = token.Method.Sign(unsignedToken, signKey); err != nil {
		return "", err
	}
	return strings.Join([]string{unsignedToken, sig}, "."), nil
}
