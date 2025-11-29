package security

import (
	"github.com/google/uuid"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/Ferrum/managers"
)

type MatrixBasedOperationControl struct {
	DataProvider *managers.DataContext
	logger       *logging.AppLogger
}

func CreateMatrixBasedOperationControl(dataProvider *managers.DataContext,
	logger *logging.AppLogger) *MatrixBasedOperationControl {
	return &MatrixBasedOperationControl{
		DataProvider: dataProvider,
		logger:       logger,
	}
}

// IsOperationAllowed checks whether operation could be performed by user or not
/* This function uses matrix based rules to control what operation could be performed by user
 *
 */
func (m *MatrixBasedOperationControl) IsOperationAllowed(objectId string, objectType data.ObjectType,
	operation OperationType, userId uuid.UUID) (bool, error) {
	// superAdmin could do everything

	return false, nil
}
