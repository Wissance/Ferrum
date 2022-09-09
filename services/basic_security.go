package services

import (
	"Ferrum/data"
	"Ferrum/dto"
	"Ferrum/managers"
)

type BasicSecurityService struct {
	DataProvider *managers.DataContext
}

func (service *BasicSecurityService) Validate(tokenIssueData *dto.TokenGenerationData) *data.OperationError {
	return &data.OperationError{}
}

func (service *BasicSecurityService) CheckCredentials(tokenIssueData *dto.TokenGenerationData) *data.OperationError {
	return &data.OperationError{}
}
