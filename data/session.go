package data

import (
	"github.com/google/uuid"
	"time"
)

type UserSession struct {
	Id              uuid.UUID
	UserId          uuid.UUID
	Started         time.Time
	Expired         time.Time
	JwtAccessToken  string
	JwtRefreshToken string
}
