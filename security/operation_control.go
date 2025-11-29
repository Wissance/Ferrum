package security

import (
	"github.com/google/uuid"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/Ferrum/managers"
)

type OperationType string

const (
	READ       OperationType = "read"
	DELETE     OperationType = "delete"
	CREATE     OperationType = "create"
	UPDATE     OperationType = "update"
	BLOCK      OperationType = "block"
	UNBLOCK    OperationType = "unblock"
	ACTIVATE   OperationType = "activate"
	DEACTIVATE OperationType = "deactivate"
)

type OperationControl interface {
	// IsOperationAllowed function that checks whether userId could be used for performing operation
	// on specified objectType identifying by objectId
	IsOperationAllowed(realmId string, objectId string, objectType data.ObjectType,
		operation OperationType, userId uuid.UUID) (bool, error)
}

func CreateOperationControlService(dataProvider *managers.DataContext,
	logger *logging.AppLogger) OperationControl {
	controlService := CreateMatrixBasedOperationControl(dataProvider, logger)
	return controlService
}
