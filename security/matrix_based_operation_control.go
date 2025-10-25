package security

import (
	"github.com/google/uuid"
	"github.com/wissance/Ferrum/data"
)

type MatrixBasedOperationControl struct {
}

// IsOperationAllowed checks whether operation could be performed by user or not
/* This function uses matrix based rules to control what operation could be performed by user
 *
 */
func (m *MatrixBasedOperationControl) IsOperationAllowed(objectId string, objectType data.ObjectType,
	operation OperationType, userId uuid.UUID) (bool, error) {
	return false, nil
}
