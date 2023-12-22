package files

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wissance/Ferrum/errors"

	"github.com/google/uuid"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/stringFormatter"
)

// FileDataManager is the simplest Data Storage without any dependencies, it uses single JSON file (it is users and clients RO auth server)
// This context type is extremely useful for simple systems
type FileDataManager struct {
	dataFile   string
	serverData data.ServerData
	logger     *logging.AppLogger
}

// CreateFileDataManagerWithInitData initializes instance of FileDataManager and sets loaded data to serverData
/* This factory function creates initialize with data instance of  FileDataManager, error reserved for usage but always nil here
 * Parameters:
 *    serverData already loaded data.ServerData from Json file in memory
 * Returns: context and error (currently is nil)
 */
func CreateFileDataManagerWithInitData(serverData *data.ServerData) (*FileDataManager, error) {
	// todo(UMV): todo provide an error handling
	mn := &FileDataManager{serverData: *serverData}
	return mn, nil
}

func CreateFileDataManager(dataFile string, logger *logging.AppLogger) (*FileDataManager, error) {
	// todo(UMV): todo provide an error handling
	mn := &FileDataManager{dataFile: dataFile, logger: logger}
	if err := mn.loadData(); err != nil {
		return nil, fmt.Errorf("loadData failed: %w", err)
	}
	return mn, nil
}

// GetRealm function for getting Realm by name
/* Searches for a realm with name realmName in serverData adn return it. Returns realm with clients but no users
 * Parameters:
 *     - realmName - name of a realm
 * Returns: Realm and error
 */
func (mn *FileDataManager) GetRealm(realmName string) (*data.Realm, error) {
	for _, e := range mn.serverData.Realms {
		// case-sensitive comparison, myapp and MyApP are different realms
		if e.Name == realmName {
			e.Users = nil
			return &e, nil
		}
	}
	return nil, errors.ErrNotFound
}

// GetUsers function for getting all Realm User
/* This function get realm by name ant extract all its users
 * Parameters:
 *     - realmName - name of a realm
 * Returns: slice of users and error
 */
func (mn *FileDataManager) GetUsers(realmName string) ([]data.User, error) {
	for _, e := range mn.serverData.Realms {
		// case-sensitive comparison, myapp and MyApP are different realms
		if e.Name == realmName {
			if len(e.Users) == 0 {
				return nil, errors.ErrZeroLength
			}
			users := make([]data.User, len(e.Users))
			for i, u := range e.Users {
				user := data.CreateUser(u)
				users[i] = user
			}
			return users, nil
		}
	}
	return nil, errors.ErrNotFound
}

// GetClient function for getting Realm Client by name
/* Searches for a client with name realmName in a realm. This function must be used after Realm was found.
 * Parameters:
 *     - realmName - realm containing clients to search
 *     - clientName - name of a client
 * Returns: Client and error
 */
func (mn *FileDataManager) GetClient(realmName string, clientName string) (*data.Client, error) {
	realm, err := mn.GetRealm(realmName)
	if err != nil {
		return nil, fmt.Errorf("GetRealm failed: %w", err)
	}

	for _, c := range realm.Clients {
		if c.Name == clientName {
			return &c, nil
		}
	}
	return nil, errors.ErrNotFound
}

// GetUser function for getting Realm User by userName
/* Searches for a user with specified name in a realm.  This function must be used after Realm was found.
 * Parameters:
 *     - realmName - realm containing users to search
 *     - userName - name of a user
 * Returns: User and error
 */
func (mn *FileDataManager) GetUser(realmName string, userName string) (data.User, error) {
	users, err := mn.GetUsers(realmName)
	if err != nil {
		return nil, fmt.Errorf("GetUsers failed: %w", err)
	}
	for _, u := range users {
		if u.GetUsername() == userName {
			return u, nil
		}
	}
	return nil, errors.ErrNotFound
}

// GetUserById function for getting Realm User by Id
/* same functions as GetUser but uses userId to search instead of username
 */
func (mn *FileDataManager) GetUserById(realmName string, userId uuid.UUID) (data.User, error) {
	users, err := mn.GetUsers(realmName)
	if err != nil {
		return nil, fmt.Errorf("GetUsers failed: %w", err)
	}
	for _, u := range users {
		if u.GetId() == userId {
			return u, nil
		}
	}
	return nil, errors.ErrNotFound
}

// CreateRealm creates new data.Realm in a data store, receive realmData unmarshalled json in a data.Realm
/*
 *
 */
func (mn *FileDataManager) CreateRealm(realmData data.Realm) error {
	return errors.ErrOperationNotSupported
}

// CreateClient creates new data.Client in a data store, requires to pass realmName (because client name is not unique), clientData is an unmarshalled json of type data.Client
func (mn *FileDataManager) CreateClient(realmName string, clientData data.Client) error {
	return errors.ErrOperationNotSupported
}

// CreateUser creates new data.User in a data store within a realm with name = realmName
func (mn *FileDataManager) CreateUser(realmName string, userData data.User) error {
	return errors.ErrOperationNotSupported
}

// UpdateRealm updates existing data.Realm in a data store within name = realmData, and new data = realmData
func (mn *FileDataManager) UpdateRealm(realmName string, realmData data.Realm) error {
	return errors.ErrOperationNotSupported
}

// UpdateClient updates existing data.Client in a data store with name = clientName and new data = clientData
func (mn *FileDataManager) UpdateClient(realmName string, clientName string, clientData data.Client) error {
	return errors.ErrOperationNotSupported
}

// UpdateUser updates existing data.User in a data store with realm name = realName, username = userName and data=userData
func (mn *FileDataManager) UpdateUser(realmName string, userName string, userData data.User) error {
	return errors.ErrOperationNotSupported
}

// DeleteRealm removes realm from data storage (Should be a CASCADE remove of all related Users and Clients)
func (mn *FileDataManager) DeleteRealm(realmName string) error {
	return errors.ErrOperationNotSupported
}

// DeleteClient removes client with name = clientName from realm with name = clientName
func (mn *FileDataManager) DeleteClient(realmName string, clientName string) error {
	return errors.ErrOperationNotSupported
}

// DeleteUser removes data.User from data store by user (userName) and realm (realmName) name respectively
func (mn *FileDataManager) DeleteUser(realmName string, userName string) error {
	return errors.ErrOperationNotSupported
}

// loadData this function loads data from JSON file (dataFile) to serverData
func (mn *FileDataManager) loadData() error {
	rawData, err := os.ReadFile(mn.dataFile)
	if err != nil {
		mn.logger.Error(stringFormatter.Format("An error occurred during config file reading: {0}", err.Error()))
		return fmt.Errorf("os.ReadFile failed: %w", err)
	}
	mn.serverData = data.ServerData{}
	if err = json.Unmarshal(rawData, &mn.serverData); err != nil {
		mn.logger.Error(stringFormatter.Format("An error occurred during data file unmarshal: {0}", err.Error()))
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	return nil
}
