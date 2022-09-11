package services

import (
	"Ferrum/data"
	"Ferrum/dto"
	"Ferrum/errors"
	"Ferrum/managers"
	"github.com/google/uuid"
)

type PasswordBasedSecurityService struct {
	DataProvider *managers.DataContext
	UserSessions map[uuid.UUID]data.UserSession
}

func Create(dataProvider *managers.DataContext) SecurityService {
	pwdSecService := &PasswordBasedSecurityService{DataProvider: dataProvider, UserSessions: map[uuid.UUID]data.UserSession{}}
	secService := SecurityService(pwdSecService)
	return secService
}

func (service *PasswordBasedSecurityService) Validate(tokenIssueData *dto.TokenGenerationData, realm *data.Realm) *data.OperationError {
	for _, c := range realm.Clients {
		if c.Name == tokenIssueData.ClientId {
			if c.Type == data.Public {
				return nil
			}

			// here we make deal with confidential client
			if c.Auth.Type == data.ClientIdAndSecrets && c.Auth.Value == tokenIssueData.ClientSecret {
				return nil
			}

		}
	}
	return &data.OperationError{Msg: errors.InvalidClientMsg, Description: errors.InvalidClientCredentialDesc}
}

func (service *PasswordBasedSecurityService) CheckCredentials(tokenIssueData *dto.TokenGenerationData, realm *data.Realm) *data.OperationError {
	user := (*service.DataProvider).GetUser(realm, tokenIssueData.Username)
	if user == nil {
		return &data.OperationError{Msg: errors.InvalidUserCredentialsMsg, Description: errors.InvalidUserCredentialsDEsc}
	}

	// todo(UMV): use hash instead raw passwords
	password := (*user).GetPassword()
	if password != tokenIssueData.Password {
		return &data.OperationError{Msg: errors.InvalidUserCredentialsMsg, Description: errors.InvalidUserCredentialsDEsc}
	}
	return nil
}

func (service *PasswordBasedSecurityService) StartOrUpdateSession(realm string, userId uuid.UUID, duration int) uuid.UUID {
	return uuid.UUID{}
}

func (service *PasswordBasedSecurityService) GetSession(realm string, userId uuid.UUID) *data.UserSession {
	return nil
}

func (service *PasswordBasedSecurityService) IsSessionExpired(realm string, userId uuid.UUID) bool {
	return false
}
