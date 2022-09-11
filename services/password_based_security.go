package services

import (
	"Ferrum/data"
	"Ferrum/dto"
	"Ferrum/errors"
	"Ferrum/managers"
	"github.com/google/uuid"
	"time"
)

type PasswordBasedSecurityService struct {
	DataProvider *managers.DataContext
	UserSessions map[string][]data.UserSession
}

func Create(dataProvider *managers.DataContext) SecurityService {
	pwdSecService := &PasswordBasedSecurityService{DataProvider: dataProvider, UserSessions: map[string][]data.UserSession{}}
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

func (service *PasswordBasedSecurityService) GetCurrentUser(realm *data.Realm, userName string) *data.User {
	return (*service.DataProvider).GetUser(realm, userName)
}

func (service *PasswordBasedSecurityService) StartOrUpdateSession(realm string, userId uuid.UUID, duration int) uuid.UUID {
	realmSessions, ok := service.UserSessions[realm]
	sessionId := uuid.New()
	// if there are no realm sessions ...
	if !ok {
		userSession := data.UserSession{Id: sessionId, UserId: userId, Started: time.Now(),
			Expired: time.Now().Add(time.Second * time.Duration(duration))}
		service.UserSessions[realm] = append(realmSessions, userSession)
		return sessionId
	}
	// realm session exists, we should find and update Expired values OR add new
	for _, s := range realmSessions {
		if s.UserId == userId {
			s.Expired = time.Now().Add(time.Second * time.Duration(duration))
			service.UserSessions[realm] = realmSessions
			return s.Id
		}
	}
	// such session does not exist, adding
	userSession := data.UserSession{Id: sessionId, UserId: userId, Started: time.Now(),
		Expired: time.Now().Add(time.Second * time.Duration(duration))}
	service.UserSessions[realm] = append(realmSessions, userSession)
	return userSession.Id
}

func (service *PasswordBasedSecurityService) GetSession(realm string, userId uuid.UUID) *data.UserSession {
	realmSessions, ok := service.UserSessions[realm]
	if !ok {
		return nil
	}
	for _, s := range realmSessions {
		if s.UserId == userId {
			return &s
		}
	}
	return nil
}

func (service *PasswordBasedSecurityService) IsSessionExpired(realm string, userId uuid.UUID) bool {
	s := service.GetSession(realm, userId)
	if s == nil {
		return true
	}
	return s.Expired.Before(time.Now())
}
