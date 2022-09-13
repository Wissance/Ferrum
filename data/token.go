package data

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type userInfo = interface{}

type JwtTokenData struct {
	IssuedAt     time.Time `json:"iat"`
	ExpiredAt    time.Time `json:"exp"`
	JwtId        uuid.UUID `json:"jti"`
	Type         string    `json:"typ"`
	Issuer       string    `json:"iss"`
	Audience     string    `json:"aud"`
	Subject      uuid.UUID `json:"sub"`
	SessionState uuid.UUID `json:"session_state"`
	SessionId    uuid.UUID `json:"sid"`
	Scope        string    `json:"scope"`
}

type TokenRefreshData struct {
	JwtTokenData
}

type AccessTokenData struct {
	JwtTokenData
	userInfo
}

func CreateAccessToken(commonData *JwtTokenData, userData *User) *AccessTokenData {
	return &AccessTokenData{JwtTokenData: *commonData, userInfo: (*userData).GetUserInfo()}
}

func (token *AccessTokenData) Valid() error {
	if token.userInfo != nil {
		return nil
	}
	return errors.New("UserInfo is null (it can't be)")
}

func CreateRefreshToken(commonData *JwtTokenData) *TokenRefreshData {
	return &TokenRefreshData{JwtTokenData: *commonData}
}

func (token *TokenRefreshData) Valid() error {
	// just formally, we don't have anything to validate, maybe in future
	return nil
}
