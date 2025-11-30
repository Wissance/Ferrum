package security

import (
	"github.com/google/uuid"
	"github.com/wissance/Ferrum/data"
	appErr "github.com/wissance/Ferrum/errors"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/Ferrum/managers"
	sf "github.com/wissance/stringFormatter"
)

type actorType string

const (
	superUser actorType = "SuperUser"
	// realmOwner is a superUser in the own realm
	realmOwner actorType = "RealmOwner"
	realmAdmin actorType = "RealmAdmin"
	realmUser  actorType = "RealmUser"
)

type matrixItem struct {
	Operations map[OperationType]bool
}

type MatrixBasedOperationControl struct {
	DataProvider *managers.DataContext
	logger       *logging.AppLogger
	matrixRules  map[actorType]matrixItem
}

//CreateMatrixBasedOperationControl a function that creates MatrixBasedOperationControl for check allowed operations
/* This function creates Operation checker based on matrix rules, in current version
 * are using only default rules that could be overwritten in future even for separate realms
 * Parameters:
 *   - dataProvider - struct implementing access to persistent objects via managers.DataContext interface
 *   - logger - logger
 * Returns: a pointer to struct MatrixBasedOperationControl
 */
//todo(UMV): allow to override matrix types, now we are using default
func CreateMatrixBasedOperationControl(dataProvider *managers.DataContext,
	logger *logging.AppLogger) *MatrixBasedOperationControl {
	return &MatrixBasedOperationControl{
		DataProvider: dataProvider,
		logger:       logger,
		matrixRules:  getDefaultMatrixRules(),
	}
}

// getDefaultMatrixRules function that constructs default set of rules based on user type, object and operation
// this function could not cover all the cases
func getDefaultMatrixRules() map[actorType]matrixItem {
	return map[actorType]matrixItem{
		superUser: {Operations: map[OperationType]bool{
			READ:       true,
			DELETE:     true,
			CREATE:     true,
			UPDATE:     true,
			BLOCK:      true,
			UNBLOCK:    true,
			ACTIVATE:   true,
			DEACTIVATE: true,
		}},
		realmOwner: {Operations: map[OperationType]bool{
			READ:       true,
			DELETE:     true,
			CREATE:     true,
			UPDATE:     true,
			BLOCK:      true,
			UNBLOCK:    true,
			ACTIVATE:   true,
			DEACTIVATE: true,
		}},
		// realmAdmin is almost realmOwner, here there is a question about DELETE operation
		// but by default it is permitted except realm delete
		realmAdmin: {Operations: map[OperationType]bool{
			READ:       true,
			DELETE:     true,
			CREATE:     true,
			UPDATE:     true,
			BLOCK:      true,
			UNBLOCK:    true,
			ACTIVATE:   true,
			DEACTIVATE: true,
		}},
		realmUser: {Operations: map[OperationType]bool{
			READ:    true,
			DELETE:  false,
			CREATE:  false,
			UPDATE:  true,
			BLOCK:   false,
			UNBLOCK: false,
			// there is a question about could user itself make account active ?
			ACTIVATE:   true,
			DEACTIVATE: false,
		}},
	}
}

// IsOperationAllowed checks whether operation could be performed by user or not
/* This function uses matrix based rules to control what operation could be performed by user
 * System Admin (settings.Admin) user that does not belong to realm could perform any operation
 * because this is a whole system admin
 */
func (m *MatrixBasedOperationControl) IsOperationAllowed(realmId string, objectId string,
	objectType data.ObjectType, operation OperationType, userId uuid.UUID) (bool, error) {
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

	var actor actorType = realmUser

	if settings != nil {
		if settings.Admin.Id == userId {
			actor = superUser
			// we should check matrix rules, however SuperUser MUST do everything
			return true, nil
		} else {
			// realm Admin could perform any operation
			// non admin could perform operation with themselves only
			realm, realReadErr := (*m.DataProvider).GetRealm(realmId)
			// can't perform anything realm doesn't exist
			if realReadErr != nil {
				wrappedErr := appErr.NewUnknownError("GetRealm", "MatrixBasedOperationControl.IsOperationAllowed", err)
				m.logger.Error(wrappedErr.Error())
				return false, wrappedErr
			}

			if realm.Owner == userId {
				actor = realmOwner
			} else {
				for _, u := range realm.Admins {
					if u == userId {
						actor = realmAdmin
						break
					}
				}
			}

			// fast check operation allowed by type
			isOperationTypeAllowed, ok := m.matrixRules[actor].Operations[operation]
			if !ok || !isOperationTypeAllowed {
				m.logger.Warn(sf.Format("Operation \"{0}\" on object of type \"{1}\" for actor type: \"{2}\" is not allowed",
					operation, objectType, actor))
				return false, nil
			}
			// some allowed operations allowed not for all objects (SPECIAL CASES):
			// 1. realmAdmin can't delete their realm, this operation could be performed by realmOwner
			if actor == realmAdmin && operation == DELETE && objectType == data.REALM {
				return false, nil
			}
			// 2. realmUser could activate themselves by link from e-mail
			if actor == realmUser && operation == ACTIVATE && objectId != userId.String() {
				return false, nil
			}

			// check user related to the realm
			_, userReadErr := (*m.DataProvider).GetUserById(realmId, userId)
			if userReadErr != nil {
				return false, userReadErr
			}

			isObjectRelatedToRealm := m.isObjectRelatedToRealm(realmId, objectId, objectType)
			return isObjectRelatedToRealm, nil
		}
	}

	return false, nil
}

func (m *MatrixBasedOperationControl) isObjectRelatedToRealm(realmId string, objectId string,
	objectType data.ObjectType) bool {
	var err error
	switch objectType {
	case data.REALM:
		_, err = (*m.DataProvider).GetRealm(objectId)
		return err == nil
	case data.CLIENT:
		_, err = (*m.DataProvider).GetClient(realmId, objectId)
		return err == nil
	case data.USER:
		objectIdUuid, idError := uuid.Parse(objectId)
		if idError != nil {
			return false
		}
		_, err = (*m.DataProvider).GetUserById(realmId, objectIdUuid)
		return err == nil
	case data.USER_FEDERATION_SERVICE_CONFIG:
		_, err = (*m.DataProvider).GetUserFederationConfig(realmId, objectId)
		return err == nil

	}
	return false
}
