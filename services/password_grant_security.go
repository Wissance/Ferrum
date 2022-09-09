package services

import (
	"Ferrum/data"
	"Ferrum/dto"
	"Ferrum/errors"
	"Ferrum/managers"
)

type PasswordBasedSecurityService struct {
	DataProvider *managers.DataContext
}

func Create(dataProvider *managers.DataContext) SecurityService {
	pwdSecService := &PasswordBasedSecurityService{DataProvider: dataProvider}
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
	return &data.OperationError{}
}
