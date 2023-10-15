package services

import (
	"github.com/google/uuid"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/dto"
)

// SecurityService is an interface that implements all checks and manipulation with sessions data
type SecurityService interface {
	// Validate checks whether provided tokenIssueData could be used for token generation or not
	Validate(tokenIssueData *dto.TokenGenerationData, realm *data.Realm) *data.OperationError
	// CheckCredentials validates provided in tokenIssueData pairs of clientId+clientSecret and username+password
	CheckCredentials(tokenIssueData *dto.TokenGenerationData, realm *data.Realm) *data.OperationError
	// GetCurrentUserByName return CurrentUser data by name
	GetCurrentUserByName(realm *data.Realm, userName string) *data.User
	// GetCurrentUserById return CurrentUser data by шв
	GetCurrentUserById(realm *data.Realm, userId uuid.UUID) *data.User
	// StartOrUpdateSession starting new session on new successful token issue request or updates existing one with new request with valid token
	StartOrUpdateSession(realm string, userId uuid.UUID, duration int, refresh int) uuid.UUID
	// AssignTokens this function creates relation between userId and issued tokens (access and refresh)
	AssignTokens(realm string, userId uuid.UUID, accessToken *string, refreshToken *string)
	// GetSession returns user session data
	GetSession(realm string, userId uuid.UUID) *data.UserSession
	// GetSessionByAccessToken returns session data by access token
	GetSessionByAccessToken(realm string, token *string) *data.UserSession
	// GetSessionByRefreshToken returns session data by access token
	GetSessionByRefreshToken(realm string, token *string) *data.UserSession
	// CheckSessionAndRefreshExpired checks is user tokens expired or not (could user use them or should get new ones)
	CheckSessionAndRefreshExpired(realm string, userId uuid.UUID) (bool, bool)
}
