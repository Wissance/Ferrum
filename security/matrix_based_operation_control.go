package security

import (
	"github.com/google/uuid"
	"github.com/wissance/Ferrum/data"
	appErr "github.com/wissance/Ferrum/errors"
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
 * System Admin (settings.Admin) user that does not belong to realm could perform any operation
 * because this is a whole system admin
 */
func (m *MatrixBasedOperationControl) IsOperationAllowed(realmId string, objectId string,
	objectType data.ObjectType,
	operation OperationType, userId uuid.UUID) (bool, error) {
	// superAdmin could do everything
	settings, err := (*m.DataProvider).GetServerSettings()
	if err != nil {
		wrappedErr := appErr.NewUnknownError("GetServerSettings", "MatrixBasedOperationControl.IsOperationAllowed", err)
		m.logger.Error(wrappedErr.Error())
		return false, wrappedErr
	}
	if settings == nil {
		// here warn only because userId could be a user from specific data.Realm
		m.logger.Warn(appErr.SecuritySettingsAreNil)
	}

	if settings != nil {
		if settings.Admin.Id == userId {
			return true, nil
		}
		// realm Admin could perform any operation
		// non admin could perform operation with themselves only
		realm, realReadErr := (*m.DataProvider).GetRealm(realmId)
		// unlike
		if realReadErr != nil {
			wrappedErr := appErr.NewUnknownError("GetRealm", "MatrixBasedOperationControl.IsOperationAllowed", err)
			m.logger.Error(wrappedErr.Error())
			return false, wrappedErr
		}

		// if userId in the realm.Admins it means that user could do everything
		for _, u := range realm.Admins {
			if u == userId {
				return true, nil
			}
		}

		// user is not admin, we should check whether it belongs to data.Realm
	}

	return false, nil
}
