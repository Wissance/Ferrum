package data

import (
	"github.com/google/uuid"
	"time"
)

type TokenRefreshData struct {
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
