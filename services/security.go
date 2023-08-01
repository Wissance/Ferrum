package services

import (
	"github.com/google/uuid"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/dto"
)

type SecurityService interface {
	Validate(tokenIssueData *dto.TokenGenerationData, realm *data.Realm) *data.OperationError
	CheckCredentials(tokenIssueData *dto.TokenGenerationData, realm *data.Realm) *data.OperationError
	GetCurrentUserByName(realm *data.Realm, userName string) *data.User
	GetCurrentUserById(realm *data.Realm, userId uuid.UUID) *data.User
	StartOrUpdateSession(realm string, userId uuid.UUID, duration int, refresh int) uuid.UUID
	AssignTokens(realm string, userId uuid.UUID, accessToken *string, refreshToken *string)
	GetSession(realm string, userId uuid.UUID) *data.UserSession
	GetSessionByAccessToken(realm string, token *string) *data.UserSession
	GetSessionByRefreshToken(realm string, token *string) *data.UserSession
	CheckSessionAndRefreshExpired(realm string, userId uuid.UUID) (bool, bool)
}
