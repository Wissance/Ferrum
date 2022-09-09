package services

import (
	"Ferrum/data"
	"Ferrum/dto"
)

type SecurityService interface {
	Validate(data *dto.TokenGenerationData) *data.OperationError
	CheckCredentials(data *dto.TokenGenerationData) *data.OperationError
}
