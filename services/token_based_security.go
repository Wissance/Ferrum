package services

import (
	"Ferrum/data"
	"Ferrum/dto"
	"Ferrum/errors"
	"Ferrum/managers"
	"github.com/google/uuid"
	"time"
)

type TokenBasedSecurityService struct {
	DataProvider *managers.DataContext
	UserSessions map[string][]data.UserSession
}

func CreateSecurityService(dataProvider *managers.DataContext) SecurityService {
	pwdSecService := &TokenBasedSecurityService{DataProvider: dataProvider, UserSessions: map[string][]data.UserSession{}}
	secService := SecurityService(pwdSecService)
	return secService
}

func (service *TokenBasedSecurityService) Validate(tokenIssueData *dto.TokenGenerationData, realm *data.Realm) *data.OperationError {
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

func (service *TokenBasedSecurityService) CheckCredentials(tokenIssueData *dto.TokenGenerationData, realm *data.Realm) *data.OperationError {
	user := (*service.DataProvider).GetUser(realm, tokenIssueData.Username)
	if user == nil {
		return &data.OperationError{Msg: errors.InvalidUserCredentialsMsg, Description: errors.InvalidUserCredentialsDesc}
	}

	// todo(UMV): use hash instead raw passwords
	password := (*user).GetPassword()
	if password != tokenIssueData.Password {
		return &data.OperationError{Msg: errors.InvalidUserCredentialsMsg, Description: errors.InvalidUserCredentialsDesc}
	}
	return nil
}

func (service *TokenBasedSecurityService) GetCurrentUser(realm *data.Realm, userName string) *data.User {
	return (*service.DataProvider).GetUser(realm, userName)
}

func (service *TokenBasedSecurityService) StartOrUpdateSession(realm string, userId uuid.UUID, duration int) uuid.UUID {
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

func (service *TokenBasedSecurityService) AssignTokens(realm string, userId uuid.UUID, accessToken *string, refreshToken *string) {
	realmSessions, ok := service.UserSessions[realm]
	if ok {
		for _, s := range realmSessions {
			if s.UserId == userId {
				s.JwtAccessToken = *accessToken
				s.JwtRefreshToken = *refreshToken
				service.UserSessions[realm] = realmSessions
			}
		}
	}
}

func (service *TokenBasedSecurityService) GetSession(realm string, userId uuid.UUID) *data.UserSession {
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

func (service *TokenBasedSecurityService) GetSessionByAccessToken(realm string, token *string) *data.UserSession {
	realmSessions, ok := service.UserSessions[realm]
	if !ok {
		return nil
	}
	for _, s := range realmSessions {
		if s.JwtAccessToken == *token {
			return &s
		}
	}
	return nil
}

func (service *TokenBasedSecurityService) IsSessionExpired(realm string, userId uuid.UUID) bool {
	s := service.GetSession(realm, userId)
	if s == nil {
		return true
	}
	return s.Expired.Before(time.Now())
}
