package services

import (
	"Ferrum/data"
	"Ferrum/dto"
	"github.com/google/uuid"
)

type SecurityService interface {
	Validate(tokenIssueData *dto.TokenGenerationData, realm *data.Realm) *data.OperationError
	CheckCredentials(tokenIssueData *dto.TokenGenerationData, realm *data.Realm) *data.OperationError
	StartOrUpdateSession(realm string, userId uuid.UUID, duration int) uuid.UUID
	GetSession(realm string, userId uuid.UUID) *data.UserSession
	IsSessionExpired(realm string, userId uuid.UUID) bool
}
