package data

import (
	"Ferrum/utils/jsontools"
	"github.com/google/uuid"
	"time"
)

type RawUserInfo interface{}

type JwtCommonInfo struct {
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
	JwtCommonInfo
}

type AccessTokenData struct {
	jwtCommonInfo JwtCommonInfo
	rawUserInfo   RawUserInfo
	ResultData    map[string]interface{}
	ResultJsonStr string
}

func CreateAccessToken(commonData *JwtCommonInfo, userData *User) *AccessTokenData {
	token := AccessTokenData{jwtCommonInfo: *commonData, rawUserInfo: (*userData).GetUserInfo()}
	token.Init()
	return &token
}

func (token *AccessTokenData) Valid() error {
	// just pass formally, we don't have anything to validate, maybe in future
	return nil
}

func CreateRefreshToken(commonData *JwtCommonInfo) *TokenRefreshData {
	return &TokenRefreshData{JwtCommonInfo: *commonData}
}

func (token *TokenRefreshData) Valid() error {
	// just pass formally, we don't have anything to validate, maybe in future
	return nil
}

func (token *AccessTokenData) Init() {
	data, str := jsontools.MergeNonIntersect[JwtCommonInfo, RawUserInfo](&token.jwtCommonInfo, &token.rawUserInfo)
	token.ResultData = data.(map[string]interface{})
	token.ResultJsonStr = str
}
